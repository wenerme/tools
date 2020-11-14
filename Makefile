build: crontimer apkindexer

ci: build

crontimer:
	go build -o dist/crontimer ./cmd/crontimer

apkindexer:
	go build -o dist/crontimer ./cmd/apkindexer


lint:
	golangci-lint run

git-hooks:
	cp ./scripts/pre-commit .git/hooks/

tidy:
	go mod tidy
