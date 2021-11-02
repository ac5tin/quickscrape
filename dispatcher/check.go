package dispatcher

import (
	"context"
	"crypto/sha512"
	"fmt"
	"quickscrape/db"

	"github.com/georgysavva/scany/pgxscan"
)

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
        SELECT url from posts
        WHERE id  = $1
    `, genPostID(url)); err != nil {
		return false
	}
	if len(*t) == 0 {
		return false
	}
	return true
}
