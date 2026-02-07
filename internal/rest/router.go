package rest

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/netbill/logium"
	"github.com/netbill/restkit/tokens"
)

type Handlers interface {
	//place class
	CreatePlaceClass(w http.ResponseWriter, r *http.Request)

	GetPlaceClass(w http.ResponseWriter, r *http.Request)
	GetPlaceClasses(w http.ResponseWriter, r *http.Request)

	OpenUpdatePlaceClassSession(w http.ResponseWriter, r *http.Request)
	ConfirmUpdatePlaceClass(w http.ResponseWriter, r *http.Request)
	DeletePlaceClassUploadIcon(w http.ResponseWriter, r *http.Request)
	CancelUpdatePlaceClass(w http.ResponseWriter, r *http.Request)

	DeletePlaceClass(w http.ResponseWriter, r *http.Request)
	ReplacePlaceClass(w http.ResponseWriter, r *http.Request)

	//place
	CreatePlace(w http.ResponseWriter, r *http.Request)

	GetPlace(w http.ResponseWriter, r *http.Request)
	GetPlaces(w http.ResponseWriter, r *http.Request)

	OpenUpdatePlaceSession(w http.ResponseWriter, r *http.Request)
	ConfirmUpdatePlace(w http.ResponseWriter, r *http.Request)
	DeletePlaceUploadIcon(w http.ResponseWriter, r *http.Request)
	DeletePlaceUploadBanner(w http.ResponseWriter, r *http.Request)
	CancelUpdatePlace(w http.ResponseWriter, r *http.Request)

	UpdatePlaceStatus(w http.ResponseWriter, r *http.Request)
	UpdatePlaceVerify(w http.ResponseWriter, r *http.Request)
}

type Middlewares interface {
	AccountAuth(
		allowedRoles ...string,
	) func(next http.Handler) http.Handler
	UpdatePlace() func(next http.Handler) http.Handler
	UpdatePlaceClass() func(next http.Handler) http.Handler
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

type Config struct {
	Port              string
	TimeoutRead       time.Duration
	TimeoutReadHeader time.Duration
	TimeoutWrite      time.Duration
	TimeoutIdle       time.Duration
}

func (s *Router) Run(ctx context.Context, cfg Config) {
	auth := s.middlewares.AccountAuth()
	sysadmin := s.middlewares.AccountAuth(
		tokens.RoleSystemAdmin,
	)
	sysmoder := s.middlewares.AccountAuth(
		tokens.RoleSystemAdmin, tokens.RoleSystemModer,
	)
	updPlace := s.middlewares.UpdatePlace()
	updPlaceClass := s.middlewares.UpdatePlaceClass()

	r := chi.NewRouter()

	r.Route("/places-svc", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/places", func(r chi.Router) {
				r.Route("/classes", func(r chi.Router) {
					r.Get("/", s.handlers.GetPlaceClasses)
					r.With(sysadmin).Post("/", s.handlers.CreatePlaceClass)
					r.With(sysadmin).Put("/replace", s.handlers.ReplacePlaceClass)

					r.Route("/{place_class_id}", func(r chi.Router) {
						r.Get("/", s.handlers.GetPlaceClass)
						r.With(sysadmin).Delete("/", s.handlers.DeletePlaceClass)

						r.With(sysadmin).Route("/update-session", func(r chi.Router) {
							r.Put("/", s.handlers.OpenUpdatePlaceClassSession)
							r.With(updPlaceClass).Delete("/", s.handlers.CancelUpdatePlaceClass)

							r.With(updPlaceClass).Put("/confirm", s.handlers.ConfirmUpdatePlaceClass)
							r.With(updPlaceClass).Delete("/upload-icon", s.handlers.DeletePlaceClassUploadIcon)
						})

					})
				})

				r.Get("/", s.handlers.GetPlaces)
				r.With(auth).Post("/", s.handlers.CreatePlace)

				r.Route("/{place_id}", func(r chi.Router) {
					r.Get("/", s.handlers.GetPlace)
					r.With(auth).Delete("/", s.handlers.DeletePlaceClassUploadIcon)
					r.With(auth).Patch("/status", s.handlers.UpdatePlaceStatus)

					r.With(sysmoder).Patch("/verify", s.handlers.UpdatePlaceVerify)

					r.With(sysadmin).Route("/update-session", func(r chi.Router) {
						r.Put("/", s.handlers.OpenUpdatePlaceClassSession)
						r.With(updPlaceClass).Delete("/", s.handlers.CancelUpdatePlaceClass)

						r.With(updPlace).Put("/confirm", s.handlers.ConfirmUpdatePlaceClass)
						r.With(updPlace).Delete("/upload-icon", s.handlers.DeletePlaceUploadIcon)
						r.With(updPlace).Delete("/upload-banner", s.handlers.DeletePlaceUploadBanner)
					})
				})
			})
		})
	})

	srv := &http.Server{
		Handler:           r,
		Addr:              cfg.Port,
		ReadTimeout:       cfg.TimeoutRead,
		ReadHeaderTimeout: cfg.TimeoutReadHeader,
		WriteTimeout:      cfg.TimeoutWrite,
		IdleTimeout:       cfg.TimeoutIdle,
	}

	s.log.Infof("starting REST service on %s", cfg.Port)

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
