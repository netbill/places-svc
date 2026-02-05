package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/netbill/logium"
	"github.com/netbill/places-svc/cmd"
	"github.com/netbill/restkit/tokens"
)

type Handlers interface {
	CreatePlaceClass(w http.ResponseWriter, r *http.Request)
	GetPlaceClass(w http.ResponseWriter, r *http.Request)
	GetPlaceClasses(w http.ResponseWriter, r *http.Request)
	UpdatePlaceClass(w http.ResponseWriter, r *http.Request)
	DeletePlaceClass(w http.ResponseWriter, r *http.Request)
	ReplacePlaceClass(w http.ResponseWriter, r *http.Request)

	CreatePlace(w http.ResponseWriter, r *http.Request)
	GetPlace(w http.ResponseWriter, r *http.Request)
	GetPlaces(w http.ResponseWriter, r *http.Request)
	UpdatePlace(w http.ResponseWriter, r *http.Request)
	UpdatePlaceStatus(w http.ResponseWriter, r *http.Request)
	UpdatePlaceVerify(w http.ResponseWriter, r *http.Request)
	DeletePlace(w http.ResponseWriter, r *http.Request)
}

type Middlewares interface {
	AccountAuth(
		allowedRoles ...string,
	) func(next http.Handler) http.Handler
	UpdatePlaces() func(next http.Handler) http.Handler
}

type Router struct {
	handlers    Handlers
	middlewares Middlewares
	log         *logium.Logger
}

func New(
	log *logium.Logger,
	middlewares Middlewares,
	handlers Handlers,
) *Router {
	return &Router{
		log:         log,
		middlewares: middlewares,
		handlers:    handlers,
	}
}

func (s *Router) Run(ctx context.Context, cfg cmd.Config) {
	auth := s.middlewares.AccountAuth()
	sysadmin := s.middlewares.AccountAuth(tokens.AccountRoleAdmin)
	sysmoder := s.middlewares.AccountAuth(tokens.AccountRoleAdmin, tokens.AccountRoleModerator)
	//updPlace := s.middlewares.UpdatePlaces()

	r := chi.NewRouter()

	r.Route("/places-svc", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/places", func(r chi.Router) {
				r.Route("/classes", func(r chi.Router) {
					r.Get("/", s.handlers.GetPlaceClasses)
					r.With(sysadmin).Post("/", s.handlers.CreatePlaceClass)

					r.Route("/{place_class_id}", func(r chi.Router) {
						r.Get("/", s.handlers.GetPlaceClass)
						r.With(sysadmin).Put("/", s.handlers.UpdatePlaceClass)
						r.With(sysadmin).Delete("/", s.handlers.DeletePlaceClass)

						r.With(sysadmin).Put("/replace", s.handlers.ReplacePlaceClass)
					})
				})

				r.Get("/", s.handlers.GetPlaces)
				r.With(auth).Post("/", s.handlers.CreatePlace)

				r.Route("/{place_id}", func(r chi.Router) {
					r.Get("/", s.handlers.GetPlace)
					r.With(auth).Put("/", s.handlers.UpdatePlace)
					r.With(auth).Delete("/", s.handlers.DeletePlace)

					r.With(auth).Patch("/status", s.handlers.UpdatePlaceStatus)
					r.With(sysmoder).Patch("/verify", s.handlers.UpdatePlaceVerify)
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
