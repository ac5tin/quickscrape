package dispatcher

func AutoDispatch() {
	t := new(tracker)
	go t.autoTrack()
}

func CrawlURL(url string) {
	queue = append(queue, url)
	go processQueue()
}
