package repository

import (
	"context"
	"errors"
	"github.com/alialin/scraperq/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type JobRepo struct {
	pool *pgxpool.Pool
}

func NewJobRepo(pool *pgxpool.Pool) *JobRepo {
	return &JobRepo{pool: pool}
}

func (j *JobRepo) Create(ctx context.Context, job *models.Job) error {

	err := j.pool.QueryRow(ctx, "INSERT INTO jobs(user_id, url, status, max_retries) VALUES ($1, $2, $3, $4) RETURNING id, created_at", job.UserID, job.URL, job.Status, job.MaxRetries).Scan(&job.ID, &job.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (j *JobRepo) FindByID(ctx context.Context, id string) (*models.Job, error) {

	var job models.Job

	err := j.pool.QueryRow(ctx, "SELECT id, user_id, url, status, retry_count, max_retries, error, created_at FROM jobs WHERE id = $1", id).
		Scan(&job.ID, &job.UserID, &job.URL, &job.Status, &job.RetryCount, &job.MaxRetries, &job.Error, &job.CreatedAt)

	if err == pgx.ErrNoRows {
		return nil, errors.New("job not found")
	}

	return &job, nil
}

func (j *JobRepo) UpdateStatus(ctx context.Context, id string, status string, jobError string) error {

	_, err := j.pool.Exec(ctx, "UPDATE jobs SET status=$1, error=$2, updated_at=NOW() WHERE id=$3", status, jobError, id)

	if err != nil {
		return err
	}

	return nil
}
