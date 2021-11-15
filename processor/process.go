package processor

import (
	"fmt"
	"log"
	"net/url"
	"quickscrape/extractor"
	"quickscrape/textprocessor"
	"sync"
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

	// assign each entity with token score of 2
	tkMap := new(sync.Map)
	wg := new(sync.WaitGroup)
	wg.Add(len(*ents))
	for _, ent := range *ents {
		go func(entity string) {
			defer wg.Done()
			lang := new(string)
			if err := tp.LangDetect(entity, lang); err != nil {
				log.Println(err.Error())
				return
			}
			entTk := new([]textprocessor.Token)
			if err := tp.Tokenise(textprocessor.InputText{Text: entity, Lang: *lang}, entTk); err != nil {
				log.Println(err.Error())
				return
			}
			for _, tk := range *entTk {
				var f float32 = 1.0
				v, _ := tkMap.LoadOrStore(tk.Token, f)
				tkMap.Store(tk.Token, v.(float32)+1)
			}
		}(ent)
	}
	wg.Wait()
	tkMap.Range(func(key, value interface{}) bool {
		*tokens = append(*tokens, textprocessor.Token{Token: key.(string), Score: value.(float32)})
		r.Tokens = append(r.Tokens, key.(string))
		return true
	})

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

	if *originSiteScore == 0 {
		*originSiteScore = 1
		if err := upsertSiteScore(&r.Site, originSiteScore); err != nil {
			return err
		}
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
	// insert
	if err := insertResults(r); err != nil {
		return err
	}
	return nil
}
