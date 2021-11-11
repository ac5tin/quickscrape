package extractor

import (
	"testing"

	"github.com/joho/godotenv"
)

func TestFullExtract(t *testing.T) {
	godotenv.Load("../.env")
	extractor := new(Extractor)
	r := new(Results)
	if err := extractor.ExtractLink("https://www.google.com", r); err != nil {
		t.Errorf("Error: %s", err)
	}

	if r.RawHTML == "" {
		t.Errorf("Failed to scrape Raw HTML")
	}

	if r.Title == "" {
		t.Errorf("Failed to scrape Title")
	}
	t.Log("success")
}
