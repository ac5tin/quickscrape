package dispatcher

import (
	"log"
	"quickscrape/extractor"
	"quickscrape/processor"
	"sync"
	"time"
)

const SCRAPE_MAX_RETRY = 5
const MAX_PARALLEL_SCRAPE = 15
const PROCESS_MAX_RETRY = 5

var queue []string = make([]string, 0)
var qchan chan string = make(chan string)

func queueProcessor() {
	queueLen := len(queue)
	log.Printf("Processing queue, Queue length: %d \n", queueLen) // debug

	scraping := 0
	wg := new(sync.WaitGroup)

	if queueLen > MAX_PARALLEL_SCRAPE {
		wg.Add(MAX_PARALLEL_SCRAPE)
	}

	for _, url := range queue {
		// scrape url
		scraping += 1
		go func(url string) error {
			if queueLen > MAX_PARALLEL_SCRAPE {
				defer wg.Done()
			}

			if checkLinkExist(url) {
				log.Printf("Link already exist: %s, skipping ...", url)
				return nil
			}

			if extractor.CheckIfLink404(url) {
				return nil
			}
			log.Printf("Scraping %s", url) // debug
			ext := new(extractor.Extractor)
			results := new(extractor.Results)
			{
				// retry scraping 3 times
				for i := 0; i < SCRAPE_MAX_RETRY; i++ {
					if err := ext.ExtractLink(url, results); err != nil {
						if i == SCRAPE_MAX_RETRY-1 {
							log.Printf("Failed to scrape %s too many times, aborting ... | ERR: %s", url, err.Error())
							return err
						}
						log.Printf("Failed to scrape %s, retrying ...", url) // debug
						continue
					}
					break
				}
			}

			log.Printf("Successfully scraped %s", url) // debug
			processor.QChan <- results

			log.Printf("Extracted %d external links + %d internal links from  %s", len(results.RelatedExternalLinks), len(results.RelatedInternalLinks), url) // debug

			for _, link := range results.RelatedExternalLinks {
				qchan <- link
			}
			for _, link := range results.RelatedInternalLinks {
				qchan <- link
			}

			return nil

		}(url)

		if queueLen > MAX_PARALLEL_SCRAPE && scraping == MAX_PARALLEL_SCRAPE {
			log.Println("max parallel scrap reached, now waiting ...")
			wg.Wait()
			wg.Add(MAX_PARALLEL_SCRAPE)
			scraping = 0
		}

	}
	queue = queue[queueLen:]
}

func processQueue() {
	go func() {
		for {
			r := <-qchan
			queue = append(queue, r)
		}
	}()
	for {
		queueProcessor()
		time.Sleep(10 * time.Second)
	}
}
