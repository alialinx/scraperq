package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/alialin/scraperq/internal/models"
	"github.com/alialin/scraperq/internal/queue"
	"github.com/alialin/scraperq/internal/worker"
)

func main() {

	q, err := queue.NewRedisQueue("localhost:6379")
	if err != nil {
		log.Fatalf("Redis bağlantı hatası: %v", err)
	}
	defer q.Close()

	log.Println("Redis bağlantısı kuruldu")

	ctx, cancel := context.WithCancel(context.Background())

	pool := worker.NewPool(3, q)

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


	go http.ListenAndServe(":8080", nil)ak
	log.Println("API Server :8080 portunda çalışıyor")

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("Kapatılıyor...")
	cancel()

	pool.Wait()

	log.Println("Temiz kapatıldı")
}

