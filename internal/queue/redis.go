package queue

import (
	"context"
	"encoding/json"
	"github.com/alialin/scraperq/internal/models"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisQueue struct {
	client *redis.Client
}

func NewRedisQueue(addr string, password string) (*RedisQueue, error) {

	// Burada oluşturduğumuz struct ile redis üzerinde yeni bir client bağlantısı açıyoruz.

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	}) // fonskiyonun aldığı adres ile redise bağlantı sağlıyoruz.

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	_, err := client.Ping(ctx).Result() // ctx ile yeni oluşruduğumuz cliente ping atıyoruz çalışıp çalışmadığını bilmek için

	if err != nil {
		return nil, err
	}

	return &RedisQueue{client: client}, nil // en son hatasız olması durumunda Yeni clienti dönderiyoruz.
}

func (q *RedisQueue) Enqueue(ctx context.Context, job *models.Job) error {

	// bu fonksiyonda redis kuyruguna yeni bir job gönderiyoruz.

	data, err := json.Marshal(job)

	if err != nil {
		return err
	}

	return q.client.LPush(ctx, "scraperq:jobs", data).Err() // LPUSH Job'ı json'a çevirip redis kuyrugunun en başına koyuyoruz.
}

func (q *RedisQueue) Dequeue(ctx context.Context, timeout time.Duration) (*models.Job, error) {
	// Bu fonksiyonda Redis kuyruk listesinin en sonundan alıyoruz.
	data, err := q.client.BRPop(ctx, timeout, "scraperq:jobs").Result()

	if err != nil {
		return nil, err
	}

	var job models.Job

	err = json.Unmarshal([]byte(data[1]), &job)
	if err != nil {
		return nil, err

	}

	return &job, nil

}

func (q *RedisQueue) Close() error {
	return q.client.Close()
}

func (q *RedisQueue) Complete(ctx context.Context, job *models.Job) error {
	data, err := json.Marshal(job)

	if err != nil {
		return err
	}

	return q.client.LRem(ctx, "scraperq:processing", 1, data).Err()

}

func (q *RedisQueue) Fail(ctx context.Context, job *models.Job) error {

	job.RetryCount++

	if job.RetryCount >= job.MaxRetries {
		job.Status = "failed"
		data, err := json.Marshal(job)
		if err != nil {
			return err
		}
		return q.client.LPush(ctx, "scraperq:dlq", data).Err()
	}

	job.Status = "pending"

	return q.Enqueue(ctx, job)

}
