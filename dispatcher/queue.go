package dispatcher

import (
	"log"
	"quickscrape/extractor"
	"quickscrape/processor"
	"time"
)

const SCRAPE_MAX_RETRY = 3

var queue []string = make([]string, 0)

func queueProcessor() {
	queueLen := len(queue)
	log.Println("Queue length:", queueLen) // debug
	for _, url := range queue {
		// scrape url
		go func(url string) error {
			if checkLinkExist(url) {
				return nil
			}
			log.Printf("Scraping %s", url) // debug
			ext := new(extractor.Extractor)
			results := new(extractor.Results)
			{
				// retry scraping 3 times
				for i := 0; i < SCRAPE_MAX_RETRY; i++ {
					if err := ext.ExtractLink(url, results); err != nil {
						log.Println(err.Error())
						if i == SCRAPE_MAX_RETRY-1 {
							return err
						}
						log.Printf("Failed to scrape %s, retrying ...", url) // debug
						continue
					}
					break
				}
			}

			log.Printf("Sending %s to indexer", url) // debug
			if err := processor.ProcessPostResults(results); err != nil {
				return err
			}

			log.Printf("Extract %d external links + %d internal links from  %s", len(results.RelatedExternalLinks), len(results.RelatedInternalLinks), url) // debug
			queue = append(queue, results.RelatedExternalLinks...)
			queue = append(queue, results.RelatedInternalLinks...)

			return nil

		}(url)
	}
	queue = queue[queueLen:]
}

func processQueue() {
	for {
		queueProcessor()
		time.Sleep(10 * time.Second)
	}
}
