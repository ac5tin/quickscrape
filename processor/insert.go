package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"quickscrape/extractor"
	"quickscrape/textprocessor"
)

type post struct {
	Author        string             `json:"author"`
	Site          string             `json:"site"`
	Title         string             `json:"title"`
	Tokens        map[string]float32 `json:"tokens"`
	Summary       string             `json:"summary"`
	URL           string             `json:"url"`
	Timestamp     uint64             `json:"timestamp"`
	Language      string             `json:"language"`
	InternalLinks []string           `json:"internal_links"`
	ExternalLinks []string           `json:"external_links"`
	Entities      map[string]float32 `json:"entities"`
}

func insertResults(r *extractor.Results) error {
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

	entTkMap := make(map[string]float32)
	for _, ent := range *ents {
		lang := new(string)
		if err := tp.LangDetect(ent, lang); err != nil {
			return err
		}

		entTk := new([]textprocessor.Token)
		if err := tp.Tokenise(textprocessor.InputText{Text: ent, Lang: *lang}, entTk); err != nil {
			return err
		}
		// assign each entity with token score of 2
		for _, tk := range *entTk {
			if v, ok := entTkMap[tk.Token]; ok {
				entTkMap[tk.Token] = v + 1
			} else {
				entTkMap[tk.Token] = 2
			}
		}
	}
	for k, v := range entTkMap {
		*tokens = append(*tokens, textprocessor.Token{Token: k, Score: v})
		r.Tokens = append(r.Tokens, k)
	}

	tkMap := make(map[string]float32)
	for _, tk := range *tokens {
		tkMap[tk.Token] = tk.Score
	}

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

	// insert post into indexer
	rr := []post{*p}
	b, err := json.Marshal(rr)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/data/insert", os.Getenv("QUICKSEARCH_ENDPOINT")), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("insert failed with status code %d", resp.StatusCode)
	}

	return nil
}
