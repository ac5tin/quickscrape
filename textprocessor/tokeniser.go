package textprocessor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Token struct {
	Token string  `json:"token"`
	Score float32 `json:"score"`
}

type InputText struct {
	Lang string `json:"lang"`
	Text string `json:"text"`
}

func (tp *TextProcessor) Tokenise(input InputText, tokens *[]Token) error {
	inp := new([]InputText)
	*inp = append(*inp, input)

	b, err := json.Marshal(inp)
	if err != nil {
		return err
	}

	for i := 0; i < MAX_RETRY_COUNT; i++ {
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/tokenise", os.Getenv("TEXTPROCESSOR_ENDPOINT")), bytes.NewBuffer(b))
		if err != nil {
			log.Printf("Failed at textprocessor.Tokenise %s", err.Error())
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed at textprocessor.Tokenise %s", err.Error())
			continue
		}
		defer resp.Body.Close()

		res := new([][]Token)
		if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
			log.Printf("Failed at textprocessor.Tokenise %s", err.Error())
			continue
		}

		if resp.StatusCode != 200 {
			log.Printf("Failed at textprocessor.Tokenise")
			continue
		}

		*tokens = (*res)[0]
		return nil

	}
	return fmt.Errorf("failed to tokenise text too many times, [lang:%s, text: %s] aborting", input.Lang, input.Text)
}

func (tp *TextProcessor) TokeniseMulti(input *[]InputText, tokens *[][]Token) error {
	inp := *input

	b, err := json.Marshal(inp)
	if err != nil {
		return err
	}

	for i := 0; i < MAX_RETRY_COUNT; i++ {
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/tokenise", os.Getenv("TEXTPROCESSOR_ENDPOINT")), bytes.NewBuffer(b))
		if err != nil {
			log.Printf("Failed at textprocessor.TokeniseMulti %s", err.Error())
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed at textprocessor.TokeniseMulti %s", err.Error())
			continue
		}
		defer resp.Body.Close()

		res := new([][]Token)
		if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
			log.Printf("Failed at textprocessor.TokeniseMulti %s", err.Error())
			continue
		}

		if resp.StatusCode != 200 {
			log.Printf("Failed at textprocessor.TokeniseMulti")
			continue
		}

		*tokens = (*res)
		return nil
	}
	return fmt.Errorf("failed to tokeniseMulti text too many times, aborting")

}
