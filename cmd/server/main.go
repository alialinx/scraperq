package main

import (
	"context"
	"encoding/json"
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
	q, err := queue.NewRedisQueue("localhost:6380")
	if err != nil {
		log.Fatalf("Redis bağlantı hatası: %v", err)
	}
	defer q.Close()

	log.Println("Redis bağlantısı kuruldu")

	ctx, cancel := context.WithCancel(context.Background())

	s := scraper.NewHTTPScraper()
	pool := worker.NewPool(3, q, s)
	pool.Start(ctx)

	log.Println("Worker pool başlatıldı (3 worker)")

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

	go http.ListenAndServe(":8080", nil)
	log.Println("API Server :8080 portunda çalışıyor")

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Kapatılıyor...")
	cancel()

	pool.Wait()

	log.Println("Temiz kapatıldı")
}
