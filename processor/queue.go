package processor

import (
	"log"
	"quickscrape/extractor"
	"sync"
	"time"
)

const MAX_PARALLEL = 3
const MAX_RETRY = 3

var Q []*extractor.Results = []*extractor.Results{}

var QChan chan *extractor.Results = make(chan *extractor.Results)

func queueProcessor(q []*extractor.Results) {
	wg := new(sync.WaitGroup)
	wg.Add(len(q))
	for _, x := range q {
		go func(r *extractor.Results) {
			defer wg.Done()
			// retry
			for i := 0; i < MAX_RETRY; i++ {
				if err := ProcessPostResults(r); err != nil {
					if i == MAX_RETRY-1 {
						log.Printf("Failed to process %s too many times, aborting ... | ERR: %s", r.URL, err.Error())
						return
					}
					log.Printf("Failed to process %s, retrying ... ", r.URL)
					continue
				}
				break
			}
			log.Printf("Successfully indexed %s", r.URL) // debug
		}(x)
	}
	wg.Wait()

}

func ProcessQueue() {
	go func() {
		for {
			r := <-QChan
			Q = append(Q, r)
		}
	}()
	for {
		capAt := len(Q)
		if len(Q) > MAX_PARALLEL {
			capAt = MAX_PARALLEL
		}
		cappedQ, newq := Q[:capAt], Q[capAt:]
		Q = newq
		queueProcessor(cappedQ)
		time.Sleep(10 * time.Second)
	}
}
