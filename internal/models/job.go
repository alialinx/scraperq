package models

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	URL        string    `json:"url"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	Result     *Result   `json:"result,omitempty"`
	Error      string    `json:"error,omitempty"`
	MaxRetries int
	RetryCount int
}

type Result struct {
	StatusCode int       `json:"status_code"`
	BodySize   int       `json:"body_size"`
	ScrapedAt  time.Time `json:"scraped_at"`
}

type ScrapeRequest struct {
	URLs []string `json:"urls"`
}

func NewJob(url string) *Job {

	id := uuid.New().String()

	status := "pending"

	createdAt := time.Now()

	return &Job{
		ID:         id,
		URL:        url,
		Status:     status,
		CreatedAt:  createdAt,
		MaxRetries: 3,
		RetryCount: 0,
	}
}
