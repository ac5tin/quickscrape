package dispatcher

import (
	"context"
	"quickscrape/db"
	"strings"

	"github.com/georgysavva/scany/pgxscan"
)

const TRACK_SITE_SCORE = 200

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
        WHERE url  = $1
        AND sites.score < $2
        AND posts.updated_at > current_date - interval '7 day'
    `, url, TRACK_SITE_SCORE); err != nil {
		return false
	}
	if len(*t) == 0 {
		return false
	}
	return true
}

// if true then remove get params from url, else leave it alone
func shouldRemoveGetParams(url string) bool {
	// keep youtube vids
	if strings.Contains(url, "youtube.com/watch?v=") {
		return false
	}
	if strings.Contains(url, "play.google.com/store/apps/details?id=") {
		return false
	}
	return true
}
