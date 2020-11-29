build: crontimer apkindexer

ci: build

crontimer:
	go build -o dist/crontimer ./cmd/crontimer

apkindexer:
	go build -o dist/crontimer ./cmd/apkindexer


lint:
	golangci-lint run

git-hooks:
	echo -e '#!/usr/bin/env bash\nmake pre-commit' > .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
	git config pull.rebase true

pre-commit:
	./scripts/pre-commit

tidy:
	go mod tidy
