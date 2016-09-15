PREFIX := /usr/local

build: get
	go build -ldflags="-s -w" -o bin/gructl -v

get:
	go get -v -t -d ./...

test:
	go test -v ./...

install: build
	install -m 0755 bin/gructl ${PREFIX}/bin/gructl

uninstall:
	rm -f ${PREFIX}/bin/gructl

clean:
	rm -f bin/gructl

format:
	go fmt .

.PHONY: build get test install uninstall clean format
