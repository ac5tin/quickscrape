package processor

import (
	"context"
	"quickscrape/db"

	"github.com/georgysavva/scany/pgxscan"
)

func getSiteScore(site *string, score *float32) error {
	conn, err := db.PG.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	scores := new([]float32)
	if err := pgxscan.Select(context.Background(), conn, scores, `
		SELECT score FROM sites WHERE site = $1
	`, *site); err != nil {
		return err
	}

	if len(*scores) == 0 {
		*score = 0
		return nil
	}

	*score = (*scores)[0]
	return nil
}

func upsertSiteScore(site *string, score *float32) error {
	conn, err := db.PG.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	if _, err := tx.Exec(context.Background(), `
		INSERT INTO sites (site, score) VALUES ($1, $2)
		ON CONFLICT (site) DO UPDATE SET score = $2
	`, *site, *score); err != nil {
		return err
	}
	return nil
}
