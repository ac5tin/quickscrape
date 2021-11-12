package processor

import (
	"fmt"
	"quickscrape/extractor"
	"quickscrape/textprocessor"
	"testing"

	"github.com/joho/godotenv"
)

func TestProcessResults(t *testing.T) {
	godotenv.Load("../.env")
	ext := new(extractor.Extractor)
	r := new(extractor.Results)
	if err := ext.ExtractLink("https://www.google.com", r); err != nil {
		t.Error(err)
	}

	// textprocessor
	tp := new(textprocessor.TextProcessor)
	// detect lang -> tokenise -> entity
	text := fmt.Sprintf("%s\n%s\n%s", r.Title, r.Summary, r.MainContent)
	if err := tp.LangDetect(text, &r.Lang); err != nil {
		t.Error(err)
	}

	tokens := new([]textprocessor.Token)
	if err := tp.Tokenise(textprocessor.InputText{Text: text, Lang: r.Lang}, tokens); err != nil {
		t.Error(err)
	}

	ents := new([]string)
	if err := tp.EntityRecognition(textprocessor.InputText{Text: text, Lang: r.Lang}, ents); err != nil {
		t.Error(err)
	}

	entTkMap := make(map[string]float32)
	for _, ent := range *ents {
		lang := new(string)
		if err := tp.LangDetect(ent, lang); err != nil {
			t.Error(err)
			return
		}

		entTk := new([]textprocessor.Token)
		if err := tp.Tokenise(textprocessor.InputText{Text: ent, Lang: *lang}, entTk); err != nil {
			t.Error(err)
			return
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

	t.Log(*tokens)
	t.Log(r.Tokens)

}
