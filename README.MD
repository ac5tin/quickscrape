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