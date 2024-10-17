#/bin/bash
# fix Failed to create table:Binary was compiled with 'CGO_ENABLED=0', go-sqlite3 requires cgo to work. This is a stub
sed -i '' 's/go 1.20/go 1.17/' go.mod && go mod tidy && CGO_ENABLED=1 GOOS=freebsd GOARCH=amd64 go build -o tgState main.go