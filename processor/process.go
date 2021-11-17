package processor

import (
	"fmt"
	"log"
	"quickscrape/extractor"
	"quickscrape/textprocessor"
	"sync"
)

const TITLE_MULTIPLIER = 2.0
const SUMMARY_MULTIPLIER = 1.5

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
	allTokens := new([][]textprocessor.Token)
	tokens := new([]textprocessor.Token)
	titleTokens := new([]textprocessor.Token)
	summaryTokens := new([]textprocessor.Token)

	everything := textprocessor.InputText{Text: text, Lang: r.Lang}
	title := textprocessor.InputText{Text: r.Title, Lang: r.Lang}
	summary := textprocessor.InputText{Text: r.Summary, Lang: r.Lang}
	if err := tp.TokeniseMulti(&[]textprocessor.InputText{everything, title, summary}, allTokens); err != nil {
		return err
	}

	*tokens = (*allTokens)[0]
	*titleTokens = (*allTokens)[1]
	*summaryTokens = (*allTokens)[2]

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

	// extra score for title
	for _, tk := range *titleTokens {
		tkMap[tk.Token] += tk.Score * TITLE_MULTIPLIER
	}

	// little extra for summary
	for _, tk := range *summaryTokens {
		tkMap[tk.Token] += tk.Score * SUMMARY_MULTIPLIER
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
	r.RelatedExternalLinks = *links
	links = nil //gc

	// - dedupe internal links
	links = new([]string)
	mapLinks = make(map[string]interface{})
	for _, l := range r.RelatedInternalLinks {
		if mapLinks[l] == nil {
			*links = append(*links, l)
			mapLinks[l] = struct{}{}
		}
	}
	mapLinks = nil // gc
	r.RelatedInternalLinks = *links
	links = nil //gc
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
