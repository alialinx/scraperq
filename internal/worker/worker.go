package worker

import (
	"context"
	"github.com/alialin/scraperq/internal/queue"
	"github.com/alialin/scraperq/internal/scraper"
	"log"
	"sync"
	"time"
)

type Pool struct {
	size    int
	queue   *queue.RedisQueue
	scraper scraper.Scraper
	wg      sync.WaitGroup
}

func NewPool(size int, q *queue.RedisQueue, s scraper.Scraper) *Pool {

	return &Pool{
		size:    size,
		queue:   q,
		scraper: s,
	}

}

func (p *Pool) Start(ctx context.Context) {

	for i := 1; i <= p.size; i++ {
		p.wg.Add(1)
		go p.runWorker(ctx, i)
	}

}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) runWorker(ctx context.Context, id int) {

	defer p.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		default:

			job, err := p.queue.Dequeue(ctx, 2*time.Second)

			if err != nil {
				continue
			}
			log.Printf("[Worker %d] İşleniyor: %s", id, job.URL)

			result, err := p.scraper.Scrape(job.URL)

			if err != nil {
				log.Printf("[Worker %d] HATA: %s → %v", id, job.URL, err)
				job.Error = err.Error()
				p.queue.Fail(ctx, job)
				continue
			}
			job.Result = result
			job.Status = "completed"
			p.queue.Complete(ctx, job)

			log.Printf("[Worker %d] Tamamlandı: %s (status: %d, size: %d)",
				id, job.URL, result.StatusCode, result.BodySize)

		}
	}

}
