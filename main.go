package main

import (
	"flag"
	"log"
	"os"
	"quickscrape/db"
	"quickscrape/dispatcher"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	site := flag.String("site", "", "URL to init auto crawl")
	autodispatch := flag.Bool("auto", false, "Auto dispatch")
	external := flag.Bool("external", true, "Enable related external links")
	flag.Parse()

	// init db
	pg, err := db.Db(os.Getenv("DB_STRING"), os.Getenv("DB_SCHEMA"))
	if err != nil {
		log.Panic(err.Error())
	}
	db.PG = pg

	// init queue
	go dispatcher.ProcessQueue()

	// dispatcher flags
	{
		// external links
		log.Println("Enable External links:", *external)
		dispatcher.EXTERNAL = *external
	}
	// -- optional dispatch --
	if *autodispatch {
		log.Println("Enabled auto dispatch")
		go dispatcher.AutoDispatch()
	}

	if *site != "" {
		log.Printf("Entry url %s", *site)
		go dispatcher.CrawlURL(*site)
	}

	for {
		// dont end program
		time.Sleep(24 * time.Hour)
	}
}
