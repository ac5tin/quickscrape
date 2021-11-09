package processor

import (
	"net/url"
	"quickscrape/extractor"
)

func ProcessPostResults(r *extractor.Results) error {
	// increment site score for all external links
	// - dedupe external links
	links := new([]string)
	mapLinks := make(map[string]interface{})
	for _, l := range r.RelatedExternalLinks {
		if mapLinks[l] == nil {
			*links = append(*links, l)
			mapLinks[l] = struct{}{}
		}
	}
	mapLinks = nil // gc
	// - each external link
	// -- externLink.site.score += r.site.score * 0.1 // max cap = 10
	// if already exist, overwrite token scores
	// insert record into to db

	originSiteScore := new(float32)
	if err := getSiteScore(&r.Site, originSiteScore); err != nil {
		return err
	}

	for _, l := range *links {
		// get site
		u, err := url.Parse(l)
		if err != nil {
			return err
		}
		// get site score
		score := new(float32)
		if err := getSiteScore(&u.Host, score); err != nil {
			return err
		}
		*score += *originSiteScore
		if err := upsertSiteScore(&u.Host, score); err != nil {
			return err
		}
	}
	return nil
}
