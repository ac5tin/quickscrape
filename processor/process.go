package processor

import "quickscrape/extractor"

func ProcessPostResults(r *extractor.Results) error {
	// increment site score for all external links
	// - dedupe external links
	links := new([]string)
	mapLinks := make(map[string]interface{})
	for _, l := range r.RelatedExternalLinks {
		if mapLinks[l] == nil {
			*links = append(*links, l)
			mapLinks[l] = nil
		}
	}
	// - each external link
	// -- externLink.site.score += r.site.score * 0.1 // max cap = 10
	// if already exist, overwrite token scores
	// insert record into to db
	return nil
}
