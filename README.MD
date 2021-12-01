# QuicKScrape
Web crawler for quicksearch



## Setup
#### Requirements
- [Spyder](https://github.com/ac5tin/spyderng)
- [TextProcessor](https://github.com/ac5tin/textprocessor)
- Postgres DB
- [QuickSearch](https://github.com/ac5tin/quicksearch)
#### Env variables
```sh
# .env
SPYDER_ENDPOINT=http://192.168.1.54:7000
TEXTPROCESSOR_ENDPOINT=http://192.168.1.54:8000

DB_STRING=postgres://
DB_SCHEMA=myschema

QUICKSEARCH_ENDPOINT=http://localhost:7898
```


## Usage
#### Flags
```txt
- site: site to crawl
- auto: auto dispatch based on site score
- external: scrape external related links
- max: max scrape per site before cooldown
- cd: site cooldown once max scrape of site is reached
- related: scrape related links
- depth: depth of the crawl
```
Example:
```sh
# only scrape one page
quickscrape -site "https://example.com" -related=false

# only scrape pages on the same host (disregard external links)
quickscrape -site "https://example.com" -external=false -cd 5 -max 15

# scrape pages and related links with each site having a cooldown of 10 minutes after 10 scrapes
quickscrape -site "https://example.com" -cd 10 -max 10
```