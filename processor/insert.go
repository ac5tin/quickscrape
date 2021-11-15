package processor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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

func insertResults(p *post) error {

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
