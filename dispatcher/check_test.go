package dispatcher

import (
	"os"
	"quickscrape/db"
	"testing"

	"github.com/joho/godotenv"
)

func TestURLExistChecker(t *testing.T) {
	godotenv.Load("../.env")
	// init db
	pg, err := db.Db(os.Getenv("DB_STRING"), os.Getenv("DB_SCHEMA"))
	if err != nil {
		panic(err)
	}
	db.PG = pg

	urls := []string{
		"https://www.yahoo.com.hk",
		"https://www.bbc.com",
		"https://www.bbc.co.uk",
	}

	for _, url := range urls {
		if !checkLinkExist(url) {
			t.Errorf("URL %s should exist", url)
		}
	}
}
