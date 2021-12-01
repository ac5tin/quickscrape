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
	sitecooldown := flag.Uint("cd", 5, "Site cooldown in minutes")
	sitemaxscrape := flag.Uint("max", 15, "Max scrape per site")
	related := flag.Bool("related", true, "Scrape related links")
	depth := flag.Int("depth", -1, "Max depth")
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
		log.Printf("Site cooldown set to %d minutes", *sitecooldown)
		dispatcher.SITE_COOLDOWN_MINUTES = *sitecooldown
		log.Printf("Max scrape per site set to %d", *sitemaxscrape)
		dispatcher.MAX_SCRAPE_PER_SITE = *sitemaxscrape
		log.Printf("Scrape related links: %t", *related)
		dispatcher.RELATED = *related
		if *depth > 0 {
			log.Printf("Max depth set to %d", *depth)
			dispatcher.DEPTH = int8(*depth)
		} else {
			log.Println("Unlimited depth")
		}
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
