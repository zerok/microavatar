all: bin/microavatar

bin:
	mkdir -p bin

bin/microavatar: $(shell find . -name '*.go') go.mod bin
	cd cmd/microavatar && go build -o ../../$@

test:
	go test ./... -v

clean:
	rm -rf bin
