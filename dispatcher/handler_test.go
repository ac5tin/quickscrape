package dispatcher

import (
	"os"
	"quickscrape/db"
	"testing"

	"github.com/joho/godotenv"
)

func TestLinkHandler(t *testing.T) {
	godotenv.Load("../.env")
	// init db
	pg, err := db.Db(os.Getenv("DB_STRING"), os.Getenv("DB_SCHEMA"))
	if err != nil {
		panic(err)
	}
	db.PG = pg

	link := "https://www.google.co.uk"
	if err := linkHandler(&link); err != nil {
		t.Error(err)
	}
}
