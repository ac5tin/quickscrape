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
	// ======== HANDLE SITE TOKEN ========
	// textprocessor
	tp := new(textprocessor.TextProcessor)
	// detect lang -> tokenise -> entity
	text := fmt.Sprintf("%s\n%s\n%s", r.Title, r.Summary, r.MainContent)
	if err := tp.LangDetect(text, &r.Lang); err != nil {
		return err
	}

	// tokenise web page
	tokens := new([]textprocessor.Token)
	if err := tp.Tokenise(textprocessor.InputText{Text: text, Lang: r.Lang}, tokens); err != nil {
		return err
	}

	// entities
	ents := new([]string)
	if err := tp.EntityRecognition(textprocessor.InputText{Text: text, Lang: r.Lang}, ents); err != nil {
		return err
	}

	// token maps
	tkMap := make(map[string]float32)    // all tokens (tokens + entities)
	entTkMap := make(map[string]float32) // entities only
	// assign token to maps
	for _, tk := range *tokens {
		tkMap[tk.Token] = tk.Score
	}
	// assign each entity with token score of 2
	entTkMapSync := new(sync.Map)
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
				v, _ := entTkMapSync.LoadOrStore(tk.Token, f)
				entTkMapSync.Store(tk.Token, v.(float32)+1)
			}
		}(ent)
	}
	wg.Wait()
	entTkMapSync.Range(func(key, value interface{}) bool {
		// add entity token to tokens map
		if v, ok := tkMap[key.(string)]; ok {
			tkMap[key.(string)] = v + value.(float32)
		} else {
			tkMap[key.(string)] = value.(float32)
		}
		// entity token map
		entTkMap[key.(string)] = value.(float32)
		return true
	})

	// +1 to all tokens (scraped tokens)
	for _, tk := range r.Tokens {
		if v, ok := tkMap[tk]; ok {
			tkMap[tk] = v + 1
		} else {
			tkMap[tk] = 1
		}
	}

	// ========= SITE ===============

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
	// ========= INSERT ==========
	// insert post
	p := new(post)
	p.Author = r.Author
	p.Site = r.Site
	p.Title = r.Title
	p.Tokens = tkMap
	p.Summary = r.Summary
	p.URL = r.URL
	p.Timestamp = r.Timestamp
	p.Language = r.Lang
	p.InternalLinks = r.RelatedInternalLinks
	p.ExternalLinks = r.RelatedExternalLinks
	p.Entities = entTkMap
	if err := insertResults(p); err != nil {
		return err
	}
	return nil
}
