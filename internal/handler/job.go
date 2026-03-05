package handler

import (
	"encoding/json"
	"net/http"

	"github.com/alialin/scraperq/internal/models"
	"github.com/alialin/scraperq/internal/queue"
	"github.com/alialin/scraperq/internal/repository"
)

type JobHandler struct {
	jobRepo *repository.JobRepo
	queue   *queue.RedisQueue
}

func NewJobHandler(jobRepo *repository.JobRepo, q *queue.RedisQueue) *JobHandler {
	return &JobHandler{jobRepo: jobRepo, queue: q}
}

func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {

	var req models.ScrapeRequest

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	jobIDs := []string{}

	for _, url := range req.URLs {

		job := models.NewJob(url)

		h.jobRepo.Create(r.Context(), job)
		h.queue.Enqueue(r.Context(), job)

		jobIDs = append(jobIDs, job.ID)

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "jobs created",
		"job_ids": jobIDs,
	})

}

func (h *JobHandler) GetJobStatus(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	data, err := h.jobRepo.FindByID(r.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "job status",
		"data":    data,
	})
}
