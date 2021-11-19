package dispatcher

import (
	"quickscrape/extractor"
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
	if err := insertResults(&[]extractor.Results{*results}); err != nil {
		return err
	}

	return nil
}
