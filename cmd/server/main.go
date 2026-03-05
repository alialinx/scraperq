package main

import (
	"context"
	"encoding/json"
	"github.com/alialin/scraperq/internal/config"
	"github.com/alialin/scraperq/internal/database"
	"github.com/alialin/scraperq/internal/models"
	"github.com/alialin/scraperq/internal/queue"
	"github.com/alialin/scraperq/internal/scraper"
	"github.com/alialin/scraperq/internal/worker"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := config.Load()
	log.Printf("Config yüklendi: %d worker, Redis: %s, Port: %s",
		cfg.WorkerCount, cfg.RedisAddr, cfg.ServerPort)

	q, err := queue.NewRedisQueue(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		log.Fatalf("Redis bağlantı hatası: %v", err)
	}
	defer q.Close()
	log.Println("Redis bağlantısı kuruldu")

	db, err := database.NewDB(cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	if err != nil {
		log.Fatalf("database error: %v", err)
	}
	defer db.Close()

	log.Printf("Database connection: %v", db)

	ctx, cancel := context.WithCancel(context.Background())

	s := scraper.NewHTTPScraper()
	pool := worker.NewPool(cfg.WorkerCount, q, s)
	pool.Start(ctx)
	log.Printf("Worker pool başlatıldı (%d worker)", cfg.WorkerCount)

	http.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "sadece POST", http.StatusMethodNotAllowed)
			return
		}

		var req models.ScrapeRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "geçersiz JSON", http.StatusBadRequest)
			return
		}

		for _, url := range req.URLs {
			job := models.NewJob(url)
			q.Enqueue(r.Context(), job)
			log.Printf("Kuyruğa eklendi: %s", url)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "jobs kuyruğa eklendi",
		})
	})

	go http.ListenAndServe(":"+cfg.ServerPort, nil)
	log.Printf("API Server :%s portunda çalışıyor", cfg.ServerPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Kapatılıyor...")
	cancel()
	pool.Wait()
	log.Println("Temiz kapatıldı")
}
