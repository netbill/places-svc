package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/netbill/places-svc/internal/media"
	"github.com/netbill/places-svc/pkg/log"
	"github.com/netbill/restkit/tokens"
)

type PlaceController interface {
	Create(w http.ResponseWriter, r *http.Request)

	Get(w http.ResponseWriter, r *http.Request)
	GetList(w http.ResponseWriter, r *http.Request)

	Update(w http.ResponseWriter, r *http.Request)

	Verify(w http.ResponseWriter, r *http.Request)
	Unverify(w http.ResponseWriter, r *http.Request)

	Activate(w http.ResponseWriter, r *http.Request)
	Deactivate(w http.ResponseWriter, r *http.Request)

	Delete(w http.ResponseWriter, r *http.Request)

	CreateUploadMediaLink(w http.ResponseWriter, r *http.Request)
	DeleteUploadMedia(w http.ResponseWriter, r *http.Request)
}

type PlaceClassController interface {
	Create(w http.ResponseWriter, r *http.Request)

	Get(w http.ResponseWriter, r *http.Request)
	GetList(w http.ResponseWriter, r *http.Request)

	Update(w http.ResponseWriter, r *http.Request)

	Deprecate(w http.ResponseWriter, r *http.Request)
	Undeprecate(w http.ResponseWriter, r *http.Request)

	CreateUploadMediaLink(w http.ResponseWriter, r *http.Request)
	DeleteUploadMedia(w http.ResponseWriter, r *http.Request)
}

type Middlewares interface {
	AccountAuth(allowedRoles ...string) func(next http.Handler) http.Handler

	Logger(log *log.Logger) func(next http.Handler) http.Handler

	CorsDocs() func(next http.Handler) http.Handler
	ResolverUrl(resolver *media.Resolver) func(next http.Handler) http.Handler
}

type Server struct {
	middlewares Middlewares

	place PlaceController
	class PlaceClassController

	log           *log.Logger
	mediaResolver *media.Resolver
}

type ServerDeps struct {
	Middlewares Middlewares

	Place PlaceController
	Class PlaceClassController

	Log           *log.Logger
	MediaResolver *media.Resolver
}

func NewServer(deps ServerDeps) *Server {
	return &Server{
		middlewares:   deps.Middlewares,
		place:         deps.Place,
		class:         deps.Class,
		log:           deps.Log,
		mediaResolver: deps.MediaResolver,
	}
}

type Config struct {
	Port              int
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
}

func (s *Server) Run(ctx context.Context, config Config) {
	auth := s.middlewares.AccountAuth()
	sysadmin := s.middlewares.AccountAuth(tokens.RoleSystemAdmin)
	sysmoder := s.middlewares.AccountAuth(tokens.RoleSystemAdmin, tokens.RoleSystemModer)

	r := chi.NewRouter()
	r.Use(
		s.middlewares.Logger(s.log),
		s.middlewares.ResolverUrl(s.mediaResolver),
		s.middlewares.CorsDocs(),
	)

	r.Route("/places-svc", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/places", func(r chi.Router) {
				r.Get("/", s.place.GetList)
				r.With(auth).Post("/", s.place.Create)

				r.Route("/{place_id}", func(r chi.Router) {
					r.Get("/", s.place.Get)
					r.With(auth).Patch("/", s.place.Update)
					r.With(auth).Delete("/", s.place.Delete)

					r.With(auth).Patch("/activate", s.place.Activate)
					r.With(auth).Patch("/deactivate", s.place.Deactivate)

					r.With(sysmoder).Route("/verify", func(r chi.Router) {
						r.Patch("/", s.place.Verify)
						r.Delete("/", s.place.Unverify)
					})

					r.With(auth).Route("/media", func(r chi.Router) {
						r.Post("/", s.place.CreateUploadMediaLink)
						r.Delete("/", s.place.DeleteUploadMedia)
					})
				})

				r.Route("/classes", func(r chi.Router) {
					r.Get("/", s.class.GetList)
					r.With(sysadmin).Post("/", s.class.Create)

					r.Route("/{place_class_id}", func(r chi.Router) {
						r.Get("/", s.class.Get)
						r.With(sysmoder).Patch("/", s.class.Update)

						r.With(sysmoder).Route("/deprecate", func(r chi.Router) {
							r.Patch("/", s.class.Deprecate)
							r.Delete("/", s.class.Undeprecate)
						})

						r.With(sysadmin).Route("/media", func(r chi.Router) {
							r.Post("/", s.class.CreateUploadMediaLink)
							r.Delete("/", s.class.DeleteUploadMedia)
						})
					})
				})
			})
		})
	})

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", config.Port),
		Handler:           r,
		ReadTimeout:       config.ReadTimeout,
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		WriteTimeout:      config.WriteTimeout,
		IdleTimeout:       config.IdleTimeout,
	}

	s.log.WithField("port", config.Port).Info("starting http service...")

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
		s.log.Info("shutting down http service...")
	case err := <-errCh:
		if err != nil {
			s.log.WithError(err).Error("http server error")
		}
	}

	shCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shCtx); err != nil {
		s.log.WithError(err).Error("failed to shutdown http server gracefully")
	} else {
		s.log.Info("http server shutdown gracefully")
	}
}
