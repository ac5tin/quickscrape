package extractor

import "testing"

func TestLinkBlockedChecker(t *testing.T) {
	link := "https://www.google.com/404"
	if blocked, code := CheckIfLinkBlocked("https://www.google.com/404"); !blocked {
		t.Errorf("Link %s should be 404, got %d", link, code)
	}

	link = "https://www.bbc.co.uk"
	if blocked, code := CheckIfLinkBlocked(link); blocked {
		t.Errorf("Link %s shouldn't be %d", link, code)
	}
}
