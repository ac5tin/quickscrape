package dispatcher

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"quickscrape/extractor"
)

func insertResults(r *[]extractor.Results) error {

	// insert post into indexer
	b, err := json.Marshal(r)
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
