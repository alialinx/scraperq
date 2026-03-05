package scraper

import (
	"github.com/alialin/scraperq/internal/models"
	"io"
	"net/http"
	"time"
)

type Scraper interface {
	Scrape(url string) (*models.Result, error)
}

type HTTPScraper struct {
	client *http.Client
}

func NewHTTPScraper() *HTTPScraper {
	client := &http.Client{Timeout: 10 * time.Second}
	return &HTTPScraper{client: client}
}

func (s *HTTPScraper) Scrape(url string) (*models.Result, error) {

	resp, err := s.client.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &models.Result{
		StatusCode: resp.StatusCode,
		BodySize:   len(body),
		ScrapedAt:  time.Now(),
	}, nil

}
