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
	RawHTML              string   `json:"rawHTML"`
	URL                  string   `json:"url"`
	Title                string   `json:"title"`
	Summary              string   `json:"summary"`
	Author               string   `json:"author"`
	MainContent          string   `json:"mainContent"`
	Timestamp            uint64   `json:"timestamp"`
	Site                 string   `json:"site"`
	Country              string   `json:"country"`
	Lang                 string   `json:"lang"`
	Type                 string   `json:"type"`
	RelatedInternalLinks []string `json:"relatedInternalLinks"`
	RelatedExternalLinks []string `json:"relatedExternalLinks"`
	Tokens               []string `json:"tokens"`
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
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/crawl/full", os.Getenv("SPYDER_ENDPOINT")), bytes.NewBuffer(b))
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

	type extractorResp struct {
		Data   Results `json:"data"`
		Status string  `json:"status"`
		Error  *string `json:"error"`
	}

	res := new(extractorResp)
	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("%s", *res.Error)
	}

	*r = res.Data

	return nil
}
