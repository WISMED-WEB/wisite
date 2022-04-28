# wisite-api

1. `./build.sh [release]`. with `release`, create executable package
2. access `http://localhost:1323/swagger/index.html`
3. before submit to github, do `./clean.sh all`. with `all`, remove all db data.

## NOTE: github.com/swaggo/echo-swagger @v1.3.0.  if @v1.3.1, blank swagger UI page

1. `update.sh`
2. change `github.com/swaggo/echo-swagger` back to `@v1.3.0`
3. `go get ./...`
4. original above steps
