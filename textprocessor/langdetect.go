package textprocessor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func (tp *TextProcessor) LangDetect(text string, lang *string) error {
	r := []string{text}
	b, err := json.Marshal(r)
	if err != nil {
		return err
	}

	for i := 0; i < MAX_RETRY_COUNT; i++ {
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/langdet", os.Getenv("TEXTPROCESSOR_ENDPOINT")), bytes.NewBuffer(b))
		if err != nil {
			log.Printf("Failed at textprocessor.LangDetect %s", err.Error())
			time.Sleep(time.Second * 3)
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Failed at textprocessor.LangDetect %s", err.Error())
			time.Sleep(time.Second * 3)
			continue
		}
		defer resp.Body.Close()

		res := new([]string)
		if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
			log.Printf("Failed at textprocessor.LangDetect %s", err.Error())
			time.Sleep(time.Second * 3)
			continue
		}

		if resp.StatusCode != 200 {
			log.Println("textprocessor failed to detect language")
			time.Sleep(time.Second * 3)
			continue
		}

		*lang = (*res)[0]
		return nil

	}

	return fmt.Errorf("failed to detect language too many times, [text: %s] aborting", text)
}
