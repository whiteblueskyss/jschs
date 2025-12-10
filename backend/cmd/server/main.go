package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/whiteblueskyss/jschs/backend/internal/config"
	"github.com/whiteblueskyss/jschs/backend/internal/db"
	"github.com/whiteblueskyss/jschs/backend/internal/handler"
	"github.com/whiteblueskyss/jschs/backend/internal/repo"
	"github.com/whiteblueskyss/jschs/backend/internal/service"
)

func main() {
	cfg := config.Load()
	if cfg.DatabaseURL == "" && cfg.ServerAddr == "" {
		log.Fatal("DATABASE_URL or SERVER_ADDR is not set. Export DATABASE_URL and SERVER_ADDR env vars (see .env.example).")
	}

	// connect to DB
	pool, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	log.Println("OK — connected to Postgres")
	defer pool.Close()

	// build repo, service, handler
	teacherRepo := repo.NewTeacherRepo(pool)
	teacherSvc := service.NewTeacherService(teacherRepo)
	teacherHandler := handler.NewTeacherHandler(teacherSvc)

	// router
	r := chi.NewRouter()
	// API v1 group
	r.Route("/api/v1", func(r chi.Router) {
		teacherHandler.Routes(r) // mounts POST /api/v1/teachers
	})

	// create http server
	srv := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: r,
	}

	// start server in goroutine
	go func() {
		log.Printf("http server listening on %s\n", cfg.ServerAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %v", err)
		}
	}()

	// graceful shutdown on interrupt
	// Wait for interrupt signal to gracefully shutdown.

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()
	log.Println("shutting down server (graceful)...")

	// 	When you press Ctrl+C in the terminal:
	// OS sends SIGINT (interrupt signal) to your process
	// The goroutine created by signal.NotifyContext catches it
	// That goroutine cancels the context
	// Cancelling the context closes the ctx.Done() channel
	// The <- operator unblocks and execution continues to line 36

	log.Println("shutting down server (cleaning resources)...")

	// final cleanup wait — demonstrate graceful close if needed
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = shutdownCtx
	// pool.Close() already deferred
	time.Sleep(200 * time.Millisecond)
}
