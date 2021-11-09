package extractor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Extractor struct{}

type Results struct {
	RawHTML              string
	URL                  string
	Title                string
	Summary              string
	Author               string
	Timestamp            uint64
	Site                 string
	Country              string
	Lang                 string
	Type                 string
	RelatedInternalLinks []string
	RelatedExternalLinks []string
	Tokens               []string
}

func (e *Extractor) ExtractLink(url string, r *Results) error {
	type ExtractorReq struct {
		URL string `json:"url"`
	}
	er := ExtractorReq{URL: url}
	b, err := json.Marshal(er)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/scrape/extract", os.Getenv("SPYDER_ENDPOINT")), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	type extractorResp struct {
		Data   Results `json:"data"`
		Result string  `json:"result"`
	}

	res := new(extractorResp)
	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		return err
	}

	if res.Result != "success" {
		return fmt.Errorf("extractor failed")
	}

	*r = res.Data

	return nil
}
