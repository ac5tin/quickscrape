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
	if err := pgxscan.Get(context.Background(), conn, score, `
		SELECT score FROM site WHERE site = $1
	`, *site); err != nil {
		return err
	}
	return nil
}
