package dispatcher

import (
	"context"
	"crypto/sha512"
	"fmt"
	"quickscrape/db"

	"github.com/georgysavva/scany/pgxscan"
)

const TRACK_SITE_SCORE = 200

func genPostID(url string) string {
	return fmt.Sprintf("%x", sha512.Sum512([]byte(url)))
}

func checkLinkExist(url string) bool {
	conn, err := db.PG.Acquire(context.Background())
	if err != nil {
		return false
	}
	defer conn.Release()

	t := new([]string)
	if err := pgxscan.Select(context.Background(), conn, t, `
		SELECT posts.url from posts
        LEFT JOIN sites
        ON posts.site = sites.site
        WHERE id  = $1
        AND sites.score < $2
        AND posts.updated_at > current_date - interval '7 day'
    `, genPostID(url), TRACK_SITE_SCORE); err != nil {
		return false
	}
	if len(*t) == 0 {
		return false
	}
	return true
}
