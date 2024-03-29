VERSION ?= v0.0.x
DIST ?= development
CLI = ./bin/unboxd
CREDENTIALS ?= .credentials.dev
CLIENT ?= .credentials.client
JWT ?= .credentials.jwt
FOLDERID ?= 147495046780
FILEID ?= 1200261074445
FILE ?= ./runtime/kandinsky.jpg

.DEFAULT_GOAL = debug

.PHONY: clean
.PHONY: update
.PHONY: format

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

vet:
	go vet ./...

lint:
	env GOOS=darwin  GOARCH=amd64 staticcheck ./...
	env GOOS=linux   GOARCH=amd64 staticcheck ./...
	env GOOS=windows GOARCH=amd64 staticcheck ./...

build-all: test vet lint
	mkdir -p dist/$(DIST)/linux
	mkdir -p dist/$(DIST)/darwin
	mkdir -p dist/$(DIST)/windows

	env GOOS=linux   GOARCH=amd64 go build -ldflags="-X 'main.VERSION=$(VERSION)'" -o ./dist/$(DIST)/linux    ./...
	env GOOS=darwin  GOARCH=amd64 go build -ldflags="-X 'main.VERSION=$(VERSION)'" -o ./dist/$(DIST)/darwin   ./...
	env GOOS=windows GOARCH=amd64 go build -ldflags="-X 'main.VERSION=$(VERSION)'" -o ./dist/$(DIST)/windows  ./...

release: build-all

debug: build
#	$(CLI) --debug --credentials $(CLIENT) list-folders
#	$(CLI) --debug --credentials $(CLIENT) list-folders --tags '/**'
#	$(CLI) --debug --credentials $(CLIENT) list-folders --tags --checkpoint ./runtime/.checkpoint --file "./runtime/folders.tsv" '/**'
#	cat ./runtime/folders.tsv
#	$(CLI) --debug --credentials $(CLIENT) list-files --batch-size 5
#	$(CLI) --debug --credentials $(CLIENT) list-folders --tags --batch-size 5 --delay 2.5s
#	$(CLI) --debug --credentials $(CLIENT) list-files   --tags --batch-size 5 --delay 2.5s
	$(CLI) help help


help: build
	$(CLI) help
	$(CLI) help list-folders
	$(CLI) help list-files
	$(CLI) help upload-file
	$(CLI) help delete-file
	$(CLI) help tag-file
	$(CLI) help untag-file
	$(CLI) help retag-file
	$(CLI) help list-templates
	$(CLI) help get-template
	$(CLI) help create-template
	$(CLI) help delete-template
	$(CLI) help version
	$(CLI) help help

version: build
	$(CLI) version

list-folders: build
#	$(CLI) --debug --credentials $(CREDENTIALS) list-folders
#	$(CLI) --debug --credentials $(CLIENT) list-folders
#	$(CLI) --debug --credentials $(CLIENT) list-folders '/'
#	$(CLI) --debug --credentials $(CLIENT) list-folders '/**'
#	$(CLI) --debug --credentials $(CLIENT) list-folders '/alpha'
#	$(CLI) --debug --credentials $(CLIENT) list-folders '/alpha/'
#	$(CLI) --debug --credentials $(CLIENT) list-folders '/alpha/*'
#	$(CLI) --debug --credentials $(CLIENT) list-folders '/alpha/pending'
	$(CLI) --debug --credentials $(CLIENT) list-folders --tags --file "./runtime/folders.tsv" '/**'
	cat "./runtime/folders.tsv"

list-files: build
#	$(CLI) --debug --credentials $(CREDENTIALS) list-files /alpha/pending
#	$(CLI) --debug --credentials $(CLIENT)      list-files /alpha/pending
#	$(CLI) --debug --credentials $(JWT)         list-files /alpha/pending
#	$(CLI) --debug --credentials $(CREDENTIALS) list-files --file "./runtime/files.tsv"
	$(CLI) --debug --credentials $(CREDENTIALS) list-files --tags --file "./runtime/files.tsv"

upload-file: build
	$(CLI) --debug --credentials $(CREDENTIALS) upload-file $(FILE) $(FOLDERID)

delete-file: build
	$(CLI) --debug --credentials $(CREDENTIALS) delete-file $(FILEID)

tag-file: build
	$(CLI) --debug --credentials $(CREDENTIALS) tag-file $(FILEID) 'woot'
	$(CLI) --debug --credentials $(CREDENTIALS) list-files --tags

untag-file: build
	$(CLI) --debug --credentials $(CREDENTIALS) untag-file $(FILEID) 'woot'
	$(CLI) --debug --credentials $(CREDENTIALS) list-files --tags

retag-file: build
	$(CLI) --debug --credentials $(CREDENTIALS) retag-file $(FILEID) 'woot' 'wooted'
	$(CLI) --debug --credentials $(CREDENTIALS) list-files --tags

list-templates: build
	$(CLI) --debug --credentials $(CREDENTIALS) list-templates

create-template: build
#	$(CLI) --debug --credentials $(CREDENTIALS) create-template
	$(CLI) --debug --credentials $(CREDENTIALS) create-template ./runtime/templates/QWERTY.json

delete-template: build
#	$(CLI) --debug --credentials $(CREDENTIALS) delete-template
	$(CLI) --debug --credentials $(CREDENTIALS) delete-template QWERTY

get-template: build
#	$(CLI) --debug --credentials $(CREDENTIALS) get-template
	$(CLI) --debug --credentials $(CREDENTIALS) get-template QWERTY
