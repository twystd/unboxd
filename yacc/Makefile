CLI=./bin/unboxd

.PHONY: clean

clean:
	go clean
	rm -rf bin

yacc:
	goyacc -o cmd/unboxd/unboxd.go grammar/unboxd.y

format: 
	go fmt ./...

build: format
	mkdir -p bin
	go build -o bin ./...

run: build
	$(CLI)
