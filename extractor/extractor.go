package extractor

type Extractor struct{}

type Results struct {
	RawHTML              string
	URL                  string
	Title                string
	Summary              string
	Author               string
	Timestamp            uint64
	Domain               string
	Country              string
	Lang                 string
	Type                 string
	RelatedInternalLinks []string
	RelatedExternalLinks []string
	Tokens               []string
}

func (e *Extractor) ExtractLink(url string, r *Results) error {
	return nil
}
