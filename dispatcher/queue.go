package dispatcher

import (
	"log"
	URL "net/url"
	"quickscrape/extractor"
	"sync"
	"time"
)

const SCRAPE_MAX_RETRY = 5
const MAX_PARALLEL_SCRAPE = 15
const PROCESS_MAX_RETRY = 5

// flags
var EXTERNAL = true

var queue []string = make([]string, 0)
var qchan chan string = make(chan string)

// ensure sites dont get blocked on too many requests
const MAX_SCRAPE_PER_SITE = 15  // max number of scrapes per site in 15 mins
const SITE_COOLDOWN_MINUTES = 5 // number of minutes to cooldown once site reached max scrape

var siteCount = new(sync.Map)

func queueProcessor() {
	queueLen := len(queue)
	log.Printf("Processing queue, Queue length: %d \n", queueLen) // debug

	scraping := 0
	wg := new(sync.WaitGroup)

	if queueLen > MAX_PARALLEL_SCRAPE {
		wg.Add(MAX_PARALLEL_SCRAPE)
	}

	for _, url := range queue {
		{
			// checking
			if checkLinkExist(url) {
				log.Printf("Link already exist: %s, skipping ...", url)
				continue
			}

			// check site max scrape reached
			u, err := URL.Parse(url)
			if err != nil {
				log.Printf("Failed to parse url of %s", url)
				continue
			}
			hostname := u.Host
			v, _ := siteCount.LoadOrStore(hostname, 0)
			log.Printf("Checking ... , Site: %s, Count: %d", hostname, v)
			// check if cooling down
			if v == -1 {
				log.Printf("Site: %s still cooling down ...", hostname)
				qchan <- url
				continue
			}
			log.Printf("Checking is %s is available ...", url) // debug
			if blocked, code := extractor.CheckIfLinkBlocked(url); blocked {
				log.Printf("Link is %d: %s, skipping ...", code, url)
				continue
			}
			// cool down if reached limit
			if v == MAX_SCRAPE_PER_SITE {
				qchan <- url
				go func() {
					siteCount.Store(hostname, -1)
					time.Sleep(time.Minute * SITE_COOLDOWN_MINUTES)
					siteCount.Store(hostname, 0)
				}()
				continue
			}
			siteCount.Store(hostname, v.(int)+1)

		}
		// scrape url
		scraping += 1
		go func(url string) error {
			if queueLen > MAX_PARALLEL_SCRAPE {
				defer wg.Done()
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
			if err := insertResults(&[]extractor.Results{*results}); err != nil {
				return err
			}

			log.Printf("Extracted %d external links + %d internal links from  %s", len(results.RelatedExternalLinks), len(results.RelatedInternalLinks), url) // debug

			if EXTERNAL {
				for _, link := range results.RelatedExternalLinks {
					qchan <- link
				}
			}
			for _, link := range results.RelatedInternalLinks {
				qchan <- link
			}

			return nil

		}(url)

		if queueLen > MAX_PARALLEL_SCRAPE && scraping == MAX_PARALLEL_SCRAPE {
			log.Println("max parallel scrape reached, now waiting ...")
			wg.Wait()
			wg.Add(MAX_PARALLEL_SCRAPE)
			scraping = 0
		}

	}
	queue = queue[queueLen:]
}

func ProcessQueue() {
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
