# URL Shortener

A Go 1.22 HTTP API with random short codes and redirects. Run `go run .` and send a POST request to `http://127.0.0.1:8080/shorten` with JSON containing a `url` field. Open the returned path on the same server. Run `go test ./...` for automated tests.

Links are stored in memory and disappear on restart. The server binds to localhost and is not intended for public deployment. Test execution remains unverified.
