PREFIX := /usr/local

build: get
	go build -o bin/gructl -v

get:
	go get -v -t ./...

test:
	go test -v ./...

install: build
	install -m 0755 bin/gructl ${PREFIX}/bin/gructl

uninstall:
	rm -f ${PREFIX}/bin/gructl

clean:
	rm -f bin/gructl

.PHONY: build get test install uninstall clean
