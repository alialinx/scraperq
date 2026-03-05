package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alialin/scraperq/internal/config"
	"github.com/alialin/scraperq/internal/database"
	"github.com/alialin/scraperq/internal/handler"
	"github.com/alialin/scraperq/internal/queue"
	"github.com/alialin/scraperq/internal/repository"
	"github.com/alialin/scraperq/internal/scraper"
	"github.com/alialin/scraperq/internal/worker"
)

func main() {
	cfg := config.Load()
	log.Printf("Config: %d worker, Redis: %s, Port: %s",
		cfg.WorkerCount, cfg.RedisAddr, cfg.ServerPort)

	q, err := queue.NewRedisQueue(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		log.Fatalf("Redis error: %v", err)
	}
	defer q.Close()
	log.Println("Redis connected")

	db, err := database.NewDB(cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)
	if err != nil {
		log.Fatalf("DB error: %v", err)
	}
	defer db.Close()
	log.Println("PostgreSQL connected")

	ctx, cancel := context.WithCancel(context.Background())

	userRepo := repository.NewUserRepo(db.Pool)
	jobRepo := repository.NewJobRepo(db.Pool)

	authHandler := handler.NewAuthHandler(userRepo, cfg.JWTSecret)
	jobHandler := handler.NewJobHandler(jobRepo, q)

	mux := handler.SetupRoutes(authHandler, jobHandler)

	s := scraper.NewHTTPScraper()
	pool := worker.NewPool(cfg.WorkerCount, q, s)
	pool.Start(ctx)
	log.Printf("Worker pool started (%d workers)", cfg.WorkerCount)

	go http.ListenAndServe(":"+cfg.ServerPort, mux)
	log.Printf("API Server running on :%s", cfg.ServerPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	cancel()
	pool.Wait()
	log.Println("Clean shutdown")
}
