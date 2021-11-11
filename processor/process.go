package processor

import (
	"fmt"
	"net/url"
	"quickscrape/extractor"
	"quickscrape/textprocessor"
)

func ProcessPostResults(r *extractor.Results) error {
	// textprocessor
	tp := new(textprocessor.TextProcessor)
	// detect lang -> tokenise -> entity
	text := fmt.Sprintf("%s\n%s\n%s", r.Title, r.Summary, r.MainContent)
	if err := tp.LangDetect(text, &r.Lang); err != nil {
		return err
	}

	tokens := new([]textprocessor.Token)
	if err := tp.Tokenise(textprocessor.InputText{Text: text, Lang: r.Lang}, tokens); err != nil {
		return err
	}

	ents := new([]string)
	if err := tp.EntityRecognition(textprocessor.InputText{Text: text, Lang: r.Lang}, ents); err != nil {
		return err
	}
	r.Tokens = append(r.Tokens, *ents...)

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
	if *originSiteScore > 10 {
		*originSiteScore = 10
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
		*score += *originSiteScore * 0.1
		if err := upsertSiteScore(&u.Host, score); err != nil {
			return err
		}
	}
	return nil
}
