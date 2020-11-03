build:
	go build -o dist/crontimer ./cmd/crontimer

ci: build
