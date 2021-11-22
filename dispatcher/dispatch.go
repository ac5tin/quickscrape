package dispatcher

func AutoDispatch() {
	t := new(tracker)
	go t.autoTrack()
}

func CrawlURL(url string) {
	qchan <- url
}
