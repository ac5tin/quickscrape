package dispatcher

import (
	"log"
	URL "net/url"
	"quickscrape/extractor"
	"strings"
	"sync"
	"time"
)

const SCRAPE_MAX_RETRY = 5
const MAX_PARALLEL_SCRAPE = 15
const PROCESS_MAX_RETRY = 5

// flags
var EXTERNAL = true
var RELATED = true
var DEPTH int8 = -1

var queue []string = make([]string, 0)
var qchan chan string = make(chan string)

// ensure sites dont get blocked on too many requests
var MAX_SCRAPE_PER_SITE uint = 15  // max number of scrapes per site in 15 mins
var SITE_COOLDOWN_MINUTES uint = 5 // number of minutes to cooldown once site reached max scrape

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
			// remove by extension
			{
				if strings.HasSuffix(u.Path, ".pdf") {
					continue
				}
				if strings.HasSuffix(u.Path, ".jpg") {
					continue
				}
				if strings.HasSuffix(u.Path, ".png") {
					continue
				}
				if strings.HasSuffix(u.Path, ".gif") {
					continue
				}
				if strings.HasSuffix(u.Path, ".css") {
					continue
				}
				if strings.HasSuffix(u.Path, ".js") {
					continue
				}
				if strings.HasSuffix(u.Path, ".ico") {
					continue
				}
				if strings.HasSuffix(u.Path, ".zip") {
					continue
				}
			}
			v, _ := siteCount.LoadOrStore(hostname, 0)
			//log.Printf("Site: %s count: %v", hostname, v) // debug
			// check if cooling down
			if v == -1 {
				log.Printf("Site: %s still cooling down ...", hostname)
				qchan <- url
				continue
			}
			// check availability
			log.Printf("Checking is %s is available ...", url) // debug
			if blocked, code := extractor.CheckIfLinkBlocked(url); blocked {
				log.Printf("Link is %d: %s, skipping ...", code, url)
				continue
			}
			// check if get request params are required
			{
				shortURL := strings.Split(url, "?")[0]
				if blocked, _ := extractor.CheckIfLinkBlocked(shortURL); !blocked {
					url = shortURL
				} else {
					log.Printf("%s is not available, fallback to %s", shortURL, url)
				}
			}

			// cool down if reached limit
			if uint(v.(int)) == MAX_SCRAPE_PER_SITE {
				qchan <- url
				go func() {
					siteCount.Store(hostname, -1)
					time.Sleep(time.Minute * time.Duration(SITE_COOLDOWN_MINUTES))
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

			log.Printf("Extracted %d external links + %d internal links from  %s", len(results.RelatedExternalLinks), len(results.RelatedInternalLinks), results.URL) // debug
			if !RELATED {
				return nil
			}

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
		if DEPTH > 0 {
			DEPTH--
		}
		if DEPTH == 0 {
			// reached end
			return
		}
		time.Sleep(10 * time.Second)
	}
}
