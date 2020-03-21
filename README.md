# mboard-go
mboard backend server

## requirements
* golang 1.8.x
* sqlite3

## build
* `git clone git@ssh.github.com:madsportslab/mboard-go`
* `go get github.com/eknkc/amber`
* `go get github.com/gorilla/mux`
* `go get github.com/gorilla/websocket`
* `go get github.com/mattes/go-sqlite3`
* `go get github.com/skip2/go-qrcode`
* `go build`
* `go install`

## setup

### initialize database
* `sqlite3 db/md.db < $GOPATH/src/github.com/madsportslab/mboard-go/db/migrations/0001_migration_init.up.sql`

### prepare web assets
