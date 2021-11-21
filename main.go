package main

import (
	"flag"
	"log"
	"os"
	"quickscrape/db"
	"quickscrape/dispatcher"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	site := flag.String("site", "", "URL to init auto crawl")
	autodispatch := flag.Bool("auto", false, "Auto dispatch")
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

	if *autodispatch {
		log.Println("Enabled auto dispatch")
		go dispatcher.AutoDispatch()
	}

	go dispatcher.CrawlURL(*site)

	for {
		// dont end program
	}
}
