.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/world world/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/emoji emoji/*
	env GOOS=linux go build -ldflags="-s -w" -o bin/people people/main.go

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

teardown:
	sls remove --verbose

test:
	go test ./... -cover

test-report:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
