package extractor

import "net/http"

func CheckIfLink404(link string) bool {
	req, err := http.NewRequest("HEAD", link, nil)
	if err != nil {
		return false
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == 404 || resp.StatusCode == 410 || resp.StatusCode == 502 || resp.StatusCode == 503 || resp.StatusCode == 429
}
