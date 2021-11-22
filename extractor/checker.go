package extractor

import "net/http"

func CheckIfLinkBlocked(link string) (bool, uint16) {
	req, err := http.NewRequest("HEAD", link, nil)
	if err != nil {
		return false, 200
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, 200
	}
	defer resp.Body.Close()

	return resp.StatusCode == 400 || resp.StatusCode == 404 || resp.StatusCode == 410 || resp.StatusCode == 502 || resp.StatusCode == 503 || resp.StatusCode == 429, uint16(resp.StatusCode)
}
