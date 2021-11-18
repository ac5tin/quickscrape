package extractor

import "testing"

func Test404Checker(t *testing.T) {
	link := "https://www.google.com/404"
	if !CheckIfLink404("https://www.google.com/404") {
		t.Errorf("Link %s should be 404", link)
	}

	link = "https://www.bbc.co.uk"
	if CheckIfLink404(link) {
		t.Errorf("Link %s shouldn't be 404", link)
	}
}
