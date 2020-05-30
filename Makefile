.PHONY: all
all: bin/microavatar

bin:
	mkdir -p bin

.PHONY: linux
linux:
	mkdir -p dist
	cd cmd/microavatar && GOOS=linux GOARCH=amd64 go build -o ../../dist/microavatar

bin/microavatar: $(shell find . -name '*.go') go.mod bin
	cd cmd/microavatar && go build -o ../../$@

.PHONY: test
test:
	go test ./... -v

.PHONY: clean
clean:
	rm -rf bin dist
