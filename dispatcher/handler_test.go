package dispatcher

import (
	"os"
	"quickscrape/db"
	"sync"
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

	wg := new(sync.WaitGroup)
	wg.Add(3)
	go func() {
		defer wg.Done()
		link := "https://www.fasta.ai"
		if err := linkHandler(&link); err != nil {
			t.Error(err)
		}
	}()
	go func() {
		defer wg.Done()
		link := "https://www.lenx.ai"
		if err := linkHandler(&link); err != nil {
			t.Error(err)
		}
	}()
	go func() {
		defer wg.Done()
		link := "https://www.thinkcol.com"
		if err := linkHandler(&link); err != nil {
			t.Error(err)
		}
	}()
	wg.Wait()

}
