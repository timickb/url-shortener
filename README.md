[![Docker Image Size](https://badgen.net/docker/stars/timickb/url_shortener_api?icon=docker&label=stars)](https://hub.docker.com/r/timickb/url_shortener_api/)
[![Build Status](https://app.travis-ci.com/timickb/url-shortener.svg?branch=dev)](https://app.travis-ci.com/timickb/url-shortener)
[![codecov](https://codecov.io/gh/timickb/url-shortener/branch/dev/graph/badge.svg?token=TLEXMS8EJA)](https://codecov.io/gh/timickb/url-shortener)

# URL Shortener API

## Build and run
`make` - build the API server

`./artifacts/bin/urlapi -store=[store_type] -config-source=[conf_source]` - run the API server

`-store` may be `local`, `db` or `improved`. 

* `local` server will use in-memory storage for URL shortenings.
* `db` server will use PostgreSQL to store data
* `improved` server will use both PostgreSQL and in-memory storage to improve URL restoring efficiency.

`-config-source` may be `env` or `file`. If `env` chosen, config parameters will be taken from environment variables,
if `file` - from file `config.yml`

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
