package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure"
	"github.com/vitalii-tkachuk/verification-service/internal/infrastructure/config"
	"github.com/vitalii-tkachuk/verification-service/internal/ui/handler/verification"
)

// Server represents abstraction over http.Server.
type Server struct {
	port            uint8
	shutdownTimeout time.Duration
	router          *chi.Mux
}

// NewServer create Server struct.
func NewServer(ctx context.Context, cfg config.Config, application *infrastructure.Application) (context.Context, *Server) {
	srv := &Server{
		port:            cfg.Port,
		shutdownTimeout: cfg.ShutdownTimeout,
		router:          chi.NewRouter(),
	}

	srv.registerMiddlewares()
	srv.registerRoutes(application)

	return serverContext(ctx), srv
}

// registerMiddlewares is used for chi.Router middleware configuration.
func (s *Server) registerMiddlewares() {
	s.router.Use(middleware.Recoverer)
}

// registerRoutes is used for chi.Router routes configuration.
func (s *Server) registerRoutes(application *infrastructure.Application) {
	s.router.Route("/verifications", func(r chi.Router) {
		r.Post("/", verification.CreateVerificationHandler(application))
		r.Get("/{verificationUuid}", verification.GetVerificationHandler(application))
		r.Patch("/{verificationUuid}/approve", verification.ApproveVerificationHandler(application))
		r.Patch("/{verificationUuid}/decline", verification.DeclineVerificationHandler(application))
	})
}

// Run starts http.Server and wait for context.Context signal to shutdown server gracefully.
func (s *Server) Run(ctx context.Context) error {
	log.Printf("Server is running on %d", s.port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server shut down", err)
		}
	}()

	<-ctx.Done()
	ctxShutDown, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return srv.Shutdown(ctxShutDown)
}

// serverContext call cancel shutdown context function in case os.Interrupt signal received.
func serverContext(ctx context.Context) context.Context {
	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt)

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		<-interruptChannel
		cancel()
	}()

	return ctx
}
