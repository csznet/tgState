#/bin/bash
sed -i '' 's/go 1.20/go 1.17/' go.mod && go mod tidy && CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -o tgState main.go