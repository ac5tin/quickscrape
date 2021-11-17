package textprocessor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func (tp *TextProcessor) EntityRecognition(inp InputText, entities *[]string) error {
	input := new([]InputText)
	*input = append(*input, inp)

	b, err := json.Marshal(input)
	if err != nil {
		return err
	}

	for i := 0; i < MAX_RETRY_COUNT; i++ {
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/entity", os.Getenv("TEXTPROCESSOR_ENDPOINT")), bytes.NewBuffer(b))
		if err != nil {
			log.Printf("Failed at textprocessor.EntityRecognition %s", err.Error())
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed at textprocessor.EntityRecognition %s", err.Error())
			continue
		}
		defer resp.Body.Close()

		res := new([][]string)
		if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
			log.Printf("Failed at textprocessor.EntityRecognition %s", err.Error())
			continue
		}
		if resp.StatusCode != 200 {
			log.Printf("textprocessor failed to recognise entity")
			continue
		}

		*entities = (*res)[0]
		return nil

	}
	return fmt.Errorf("failed to detect entities too many times, [lang:%s, text: %s] aborting", inp.Lang, inp.Text)
}
