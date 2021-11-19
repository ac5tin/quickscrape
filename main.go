package main

import (
	"flag"
	"log"
	"os"
	"quickscrape/db"
	"quickscrape/dispatcher"
	"quickscrape/processor"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	site := flag.String("site", "", "URL to init auto crawl")
	flag.Parse()

	if *site == "" {
		log.Panic("no site supplied")
	}
	// init db
	pg, err := db.Db(os.Getenv("DB_STRING"), os.Getenv("DB_SCHEMA"))
	if err != nil {
		log.Panic(err.Error())
	}
	db.PG = pg

	go dispatcher.AutoDispatch()
	go dispatcher.CrawlURL(*site)
	go processor.ProcessQueue()

	for {
		// dont end program
	}
}
