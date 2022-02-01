VERSION ?= v0.0.x
DIST ?= development
CLI = ./bin/unboxd
CREDENTIALS ?= .credentials.dev
CLIENT ?= .credentials.client
JWT ?= .credentials.jwt
FOLDERID ?= 147495046780
FILEID ?= 903401361197
FILE ?= ./runtime/kandinsky.jpg

.PHONY: clean

all: test      \
	 benchmark \
     coverage

clean:
	rm -rf bin/*

update:
	go get -u github.com/cristalhq/jwt/v4
	go get -u github.com/google/uuid
	go get -u github.com/youmark/pkcs8

format: 
	go fmt ./...

build: format
	mkdir -p bin
	go build -ldflags="-X 'main.VERSION=$(VERSION)'" -o bin ./...

test: build
	go test -v ./...

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
	# dlv exec ./bin/boxd-cli -- help
	$(CLI) --debug --credentials $(CLIENT) list-folders '/'
	$(CLI) --debug --credentials $(CLIENT) list-folders '/*'
	$(CLI) --debug --credentials $(CLIENT) list-folders '/**'
	$(CLI) --debug --credentials $(CLIENT) list-folders '/alpha/*'

help: build
	$(CLI) --debug help

version: build
	$(CLI) --debug version

list-folders: build
	$(CLI) --debug --credentials $(CLIENT) list-folders
	$(CLI) --debug --credentials $(CLIENT) list-folders '/'
	$(CLI) --debug --credentials $(CLIENT) list-folders '/alpha'
	$(CLI) --debug --credentials $(CLIENT) list-folders '/alpha/'
	$(CLI) --debug --credentials $(CLIENT) list-folders '/alpha/*'
	$(CLI) --debug --credentials $(CLIENT) list-folders '/alpha/pending'
	$(CLI) --debug --credentials $(CREDENTIALS) list-folders

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
	$(CLI) --debug --credentials $(CLIENT) list-files /alpha/pending
	$(CLI) --debug --credentials $(JWT) list-files /alpha/pending

upload-file: build
	$(CLI) --debug --credentials $(CREDENTIALS) upload-file $(FILE) $(FOLDERID)

delete-file: build
	$(CLI) --debug --credentials $(CREDENTIALS) delete-file $(FILEID)

tag-file: build
	$(CLI) --debug --credentials $(CREDENTIALS) tag-file 907642054572 'woot'
	$(CLI) --debug --credentials $(CREDENTIALS) list-files /alpha/pending

untag-file: build
	$(CLI) --debug --credentials $(CREDENTIALS) untag-file 907642054572 'woot'
	$(CLI) --debug --credentials $(CREDENTIALS) list-files /alpha/pending

retag-file: build
	$(CLI) --debug --credentials $(CREDENTIALS) retag-file 907642054572 'woot' 'woot2'
	$(CLI) --debug --credentials $(CREDENTIALS) list-files /alpha/pending
