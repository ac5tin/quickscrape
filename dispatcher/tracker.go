package dispatcher

import (
	"context"
	"log"
	"quickscrape/db"
	"time"

	"github.com/georgysavva/scany/pgxscan"
)

type tracker struct{}

func (t *tracker) autoTrack() {
	// every 15 mins scrape sites with score >= 20
	go func() {
		for {
			sites := new([]site)
			if err := getTopSites(20, sites); err != nil {
				log.Printf("Failed to get top sites, ERR: %s", err.Error())
			}
			for _, s := range *sites {
				queue = append(queue, s.Site)
			}
			time.Sleep(15 * time.Minute)
		}
	}()
	// every 30 mins scrape sites with score >= 15
	go func() {
		for {
			sites := new([]site)
			if err := getTopSites(15, sites); err != nil {
				log.Printf("Failed to get top sites, ERR: %s", err.Error())
			}
			for _, s := range *sites {
				if s.Score >= 20 {
					// skip to avoid duplicate scrapes
					continue
				}
				queue = append(queue, s.Site)
			}
			time.Sleep(30 * time.Minute)
		}

	}()
	// every hour scrape sites with score >= 10
	go func() {
		for {
			sites := new([]site)
			if err := getTopSites(10, sites); err != nil {
				log.Printf("Failed to get top sites, ERR: %s", err.Error())
			}
			for _, s := range *sites {
				if s.Score >= 15 {
					// skip to avoid duplicate scrapes
					continue
				}
				queue = append(queue, s.Site)
			}
			time.Sleep(1 * time.Hour)
		}

	}()
	// every day scrape sites with score >= 5
	go func() {
		for {
			sites := new([]site)
			if err := getTopSites(5, sites); err != nil {
				log.Printf("Failed to get top sites, ERR: %s", err.Error())
			}
			for _, s := range *sites {
				if s.Score >= 10 {
					// skip to avoid duplicate scrapes
					continue
				}
				queue = append(queue, s.Site)
			}
			time.Sleep(24 * time.Hour)
		}

	}()
}

type site struct {
	Site  string  `db:"site"`
	Score float32 `db:"score"`
}

func getTopSites(score float32, sites *[]site) error {
	conn, err := db.PG.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	if err := pgxscan.Select(context.Background(), conn, sites, `
        SELECT site FROM sites
        WHERE score >= $1
    `, score); err != nil {
		return err
	}
	return nil
}
