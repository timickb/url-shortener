# URL Shortener API

## Endpoints

* `POST /api/v1/create` - returns a 10 symbol hash (a-z, A-Z, 0-9 and \_) for given URL

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

* `GET /api/v1/restore?shortened={hash}` - restores an original string by given {hash}