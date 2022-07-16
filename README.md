[![Build Status](https://app.travis-ci.com/timickb/url-shortener.svg?branch=dev)](https://app.travis-ci.com/timickb/url-shortener)
[![codecov](https://codecov.io/gh/timickb/url-shortener/branch/dev/graph/badge.svg?token=TLEXMS8EJA)](https://codecov.io/gh/timickb/url-shortener)

# URL Shortener API

## Build and run
`make` - build the API server

`./artifacts/bin/urlapi -store=[store_type] -config-source=[conf_source]` - run the API server

`-store` may be `local` or `db`. If chosen `local`, server will use in-memory
storage for URL shortenings. If chosen `db`, PostgreSQL will be used.

## Run in docker

`docker pull timickb/url_shortener_api:latest` - pulls repository

`docker-compose up` runs a container with PostgreSQL which reads config parameters
from `.env` file


## Endpoints

* `POST /create` - returns a 10 symbol hash (a-z, A-Z, 0-9 and \_) for given URL

Request body example:

```json
{
    "url": "https://example.com"
}
```

Response example:
```json
{
    "hash": "c7f_a366h_"
}
```

* `GET /restore?hash=[hash]` - restores an original string by given `hash`