package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/whiteblueskyss/jschs/backend/internal/config"
	"github.com/whiteblueskyss/jschs/backend/internal/db"
)

func main() {
	// load config
	cfg := config.Load()
	if cfg.DatabaseURL == "" || cfg.ServerAddr == "" {
		log.Fatal("DATABASE_URL or SERVER_ADDR is not set. (see .env.example).")
	}

	// connect to DB
	pool, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	log.Println("OK — connected to Postgres")
	// ensure pool closed on exit
	defer pool.Close()

	// Wait for interrupt signal to gracefully shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()

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
