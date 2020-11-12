build: crontimer apkindexer

ci: build

crontimer:
	go build -o dist/crontimer ./cmd/crontimer

apkindexer:
	go build -o dist/crontimer ./cmd/apkindexer
