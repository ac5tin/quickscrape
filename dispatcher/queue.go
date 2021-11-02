package dispatcher

import (
	"log"
	"quickscrape/extractor"
	"time"
)

var queue []string = make([]string, 0)

func queueProcessor() {
	for i, url := range queue {
		// scrape url
		queue = queue[i+1:]
		go func(url string) error {
			if checkLinkExist(url) {
				return nil
			}
			ext := new(extractor.Extractor)
			results := new(extractor.Results)
			if err := ext.ExtractLink(url, results); err != nil {
				log.Println(err.Error())
				return err
			}

			queue = append(queue, results.RelatedExternalLinks...)
			queue = append(queue, results.RelatedInternalLinks...)

			return nil

		}(url)
	}
}

func processQueue() {
	for {
		queueProcessor()
		time.Sleep(10 * time.Second)
	}
}
