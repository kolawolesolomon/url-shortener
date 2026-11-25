# URL Shortener

[![Tests](https://github.com/kolawolesolomon/url-shortener/actions/workflows/tests.yml/badge.svg)](https://github.com/kolawolesolomon/url-shortener/actions/workflows/tests.yml)

A Go HTTP API that creates random short paths and redirects them to HTTP or HTTPS destinations. The server uses the standard library.

## Setup and use

Requires Go 1.22+. Start the server from the repository root.

```sh
go run .
```

In another terminal, send a URL to the API, then open the returned path on the same server.

```sh
curl -i -X POST http://127.0.0.1:8080/shorten \
  -H 'Content-Type: application/json' \
  -d '{"url":"https://example.com"}'
```

Valid requests return HTTP 201 with a `path` in JSON. A GET to that path returns a 302 redirect. Invalid URLs return 400 and unknown codes return 404.

## Build and test

```sh
go build ./...
go test -race ./...
```

Handler tests cover creation, redirects, invalid input, unknown codes and concurrent requests.

## Limitations

The server binds to `127.0.0.1:8080`. It holds up to 10,000 links in memory and loses them on restart. It is intended for local use.
