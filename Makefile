DEPENDENCIES := github.com/docopt/docopt-go \
				github.com/gorilla/mux \
				github.com/stretchr/testify


all: build test

build:
		go build statsfetcher.go 

test:
		go test -v $(PACKAGES)

format:
		go fmt $(PACKAGES)

deps:
		go get $(DEPENDENCIES)

