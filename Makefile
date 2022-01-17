VERSION ?= v0.0.x
DIST ?= development
CLI = ./bin/boxd-cli
CREDENTIALS ?= .credentials.dev
FILEID ?= 903401361197

.PHONY: clean

all: test      \
	 benchmark \
     coverage

clean:
	rm -rf bin/*

format: 
	go fmt ./...

build: format
	mkdir -p bin
	go build -ldflags="-X 'main.version=$(VERSION)'" -o bin ./...

test: build
	go test ./...

benchmark: test
	go test -bench ./...

coverage: build
	go test -cover ./...

vet: build
	go vet ./...

build-all: test vet
	mkdir -p dist/$(DIST)/linux
	mkdir -p dist/$(DIST)/darwin
	mkdir -p dist/$(DIST)/windows

	env GOOS=linux   GOARCH=amd64 go build -ldflags="-X 'main.version=$(VERSION)'" -o ./dist/$(DIST)/linux    ./...
	env GOOS=darwin  GOARCH=amd64 go build -ldflags="-X 'main.version=$(VERSION)'" -o ./dist/$(DIST)/darwin   ./...
	env GOOS=windows GOARCH=amd64 go build -ldflags="-X 'main.version=$(VERSION)'" -o ./dist/$(DIST)/windows  ./...

release: build-all

debug: build
	# dlv test github.com/uhppoted/uhppoted-httpd/system/catalog
	dlv exec ./bin/boxd-cli -- help

help: build
	$(CLI) --debug help

list-templates: build
	$(CLI) --debug --credentials $(CREDENTIALS) list-templates

create-template: build
	$(CLI) --debug --credentials $(CREDENTIALS) create-template

delete-template: build
	$(CLI) --debug --credentials $(CREDENTIALS) delete-template XXX

get-template: build
	$(CLI) --debug --credentials $(CREDENTIALS) get-template PWA

list-files: build
	$(CLI) --debug --credentials $(CREDENTIALS) list-files /alpha/pending

delete-file: build
	$(CLI) --debug --credentials $(CREDENTIALS) delete-file $(FILEID)
