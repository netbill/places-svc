package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/netbill/places-svc/pkg/log"
	"github.com/netbill/restkit/tokens"
)

type Handlers interface {
	CreatePlaceClass(w http.ResponseWriter, r *http.Request)

	GetPlaceClass(w http.ResponseWriter, r *http.Request)
	GetPlaceClasses(w http.ResponseWriter, r *http.Request)

	UpdatePlaceClass(w http.ResponseWriter, r *http.Request)
	DeprecatePlaceClass(w http.ResponseWriter, r *http.Request)

	CreatePlaceClassUploadMediaLink(w http.ResponseWriter, r *http.Request)
	DeletePlaceClassUploadIcon(w http.ResponseWriter, r *http.Request)

	CreatePlace(w http.ResponseWriter, r *http.Request)

	GetPlace(w http.ResponseWriter, r *http.Request)
	GetPlaces(w http.ResponseWriter, r *http.Request)

	UpdatePlace(w http.ResponseWriter, r *http.Request)
	UpdatePlaceStatus(w http.ResponseWriter, r *http.Request)
	UpdatePlaceVerify(w http.ResponseWriter, r *http.Request)

	CreatePlaceUploadMediaLink(w http.ResponseWriter, r *http.Request)
	DeletePlaceUploadIcon(w http.ResponseWriter, r *http.Request)
	DeletePlaceUploadBanner(w http.ResponseWriter, r *http.Request)

	DeletePlace(w http.ResponseWriter, r *http.Request)
}

type Middlewares interface {
	AccountAuth(
		allowedRoles ...string,
	) func(next http.Handler) http.Handler
	Logger(log *log.Logger) func(next http.Handler) http.Handler
	CorsDocs() func(next http.Handler) http.Handler
}

type Server struct {
	handlers    Handlers
	middlewares Middlewares
}

func New(
	middlewares Middlewares,
	handlers Handlers,
) *Server {
	return &Server{
		middlewares: middlewares,
		handlers:    handlers,
	}
}

type Config struct {
	Port              int
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

func (s *Server) Run(ctx context.Context, log *log.Logger, cfg Config) {
	auth := s.middlewares.AccountAuth()
	sysadmin := s.middlewares.AccountAuth(tokens.RoleSystemAdmin)
	sysmoder := s.middlewares.AccountAuth(tokens.RoleSystemAdmin, tokens.RoleSystemModer)

	r := chi.NewRouter()
	r.Use(
		s.middlewares.Logger(log),
		s.middlewares.CorsDocs(),
	)

	r.Route("/places-svc", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/places", func(r chi.Router) {
				r.Route("/classes", func(r chi.Router) {
					r.Get("/", s.handlers.GetPlaceClasses)
					r.With(sysadmin).Post("/", s.handlers.CreatePlaceClass)

					r.Route("/{place_class_id}", func(r chi.Router) {
						r.Get("/", s.handlers.GetPlaceClass)
						r.With(sysmoder).Put("/", s.handlers.UpdatePlaceClass)
						r.With(sysadmin).Delete("/", s.handlers.DeprecatePlaceClass)

						r.With(sysadmin).Route("/media", func(r chi.Router) {
							r.Route("/upload", func(r chi.Router) {
								r.Post("/url", s.handlers.CreatePlaceClassUploadMediaLink)

								r.Delete("/icon", s.handlers.DeletePlaceClassUploadIcon)
							})
						})
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

					r.With(auth).Route("/media", func(r chi.Router) {
						r.Route("/upload", func(r chi.Router) {
							r.Post("/url", s.handlers.CreatePlaceUploadMediaLink)

							r.Delete("/icon", s.handlers.DeletePlaceUploadIcon)
							r.Delete("/banner", s.handlers.DeletePlaceUploadBanner)
						})
					})
				})
			})
		})
	})

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           r,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	log.WithField("port", cfg.Port).Info("starting http service...")

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
		log.Info("shutting down http service...")
	case err := <-errCh:
		if err != nil {
			log.WithError(err).Error("http server error")
		}
	}

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		log.WithError(err).Error("failed to shutdown http server gracefully")
	} else {
		log.Info("http server shutdown gracefully")
	}
}
