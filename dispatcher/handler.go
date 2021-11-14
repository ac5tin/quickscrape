package dispatcher

import (
	"quickscrape/extractor"
	"quickscrape/processor"
)

func linkHandler(link *string) error {
	if checkLinkExist(*link) {
		return nil
	}
	ext := new(extractor.Extractor)
	results := new(extractor.Results)
	if err := ext.ExtractLink(*link, results); err != nil {
		return err
	}
	if err := processor.ProcessPostResults(results); err != nil {
		return err
	}
	return nil
}
