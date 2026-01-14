package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/netbill/logium"
	"github.com/netbill/places-svc/internal"
	"github.com/netbill/restkit/roles"
)

type Handlers interface {
	CreatePlaceClass(w http.ResponseWriter, r *http.Request)
	GetPlaceClass(w http.ResponseWriter, r *http.Request)
	UpdatePlaceClass(w http.ResponseWriter, r *http.Request)
	DeletePlaceClass(w http.ResponseWriter, r *http.Request)
	ReplacePlaceClass(w http.ResponseWriter, r *http.Request)
}

type Middlewares interface {
	Auth() func(http.Handler) http.Handler
	RoleGrant(allowedRoles map[string]bool) func(http.Handler) http.Handler
}

type Service struct {
	handlers    Handlers
	middlewares Middlewares
	log         logium.Logger
}

func New(
	log logium.Logger,
	middlewares Middlewares,
	handlers Handlers,
) *Service {
	return &Service{
		log:         log,
		middlewares: middlewares,
		handlers:    handlers,
	}
}

func (s *Service) Run(ctx context.Context, cfg internal.Config) {
	auth := s.middlewares.Auth()
	sysadmin := s.middlewares.RoleGrant(map[string]bool{
		roles.SystemAdmin: true,
	})
	sysmoder := s.middlewares.RoleGrant(map[string]bool{
		roles.SystemAdmin: true,
		roles.SystemModer: true,
	})

	r := chi.NewRouter()

	r.Route("/places-svc", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/places", func(r chi.Router) {
				r.Route("/classes", func(r chi.Router) {
					r.Get("/", s.handlers.GetPlaceClasses)
					r.With(auth, sysadmin).Post("/", s.handlers.CreatePlaceClass)

					r.Route("/{class_id}", func(r chi.Router) {
						r.Get("/", s.handlers.GetPlaceClass)
						r.With(auth, sysadmin).Put("/", s.handlers.UpdatePlaceClass)
						r.With(auth, sysadmin).Delete("/", s.handlers.DeletePlaceClass)

						r.With(auth, sysadmin).Put("/replace", s.handlers.ReplacePlaceClass)
					})
				})

				r.Get("/", s.handlers.GetPlaces)
				r.With(auth).Post("/", s.handlers.CreatePlaces)

				r.Route("/{place_id}", func(r chi.Router) {
					r.Get("/", s.handlers.GetPlace)
					r.With(auth).Put("/", s.handlers.UpdatePlace)
					r.With(auth).Delete("/", s.handlers.DeletePlace)

					r.With(auth).Patch("/status", s.handlers.UpdatePlaceStatus)
					r.With(sysadmin).Patch("/class", s.handlers.UpdatePlacesClass)
					r.With(auth, sysmoder).Patch("/verify", s.handlers.VerifyPlace)
				})
			})
		})
	})

	srv := &http.Server{
		Handler:           r,
		Addr:              cfg.Rest.Port,
		ReadTimeout:       cfg.Rest.Timeouts.Read,
		ReadHeaderTimeout: cfg.Rest.Timeouts.ReadHeader,
		WriteTimeout:      cfg.Rest.Timeouts.Write,
		IdleTimeout:       cfg.Rest.Timeouts.Idle,
	}

	s.log.Infof("starting REST service on %s", cfg.Rest.Port)

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		} else {
			errCh <- nil
		}
	}()

	select {
	case <-ctx.Done():
		s.log.Warnf("shutting down REST service...")
	case err := <-errCh:
		if err != nil {
			s.log.Errorf("REST server error: %v", err)
		}
	}

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		s.log.Errorf("REST shutdown error: %v", err)
	} else {
		s.log.Warnf("REST server stopped")
	}
}
