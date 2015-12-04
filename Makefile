PREFIX := /usr/local

get:
	go get -v ./...

build:
	go build -o bin/gructl -v

test:
	go test -v ./...

integration:
	go test -v --tags integration ./...

install: build
	install -m 0755 bin/gructl ${PREFIX}/bin/gructl

uninstall:
	rm -f ${PREFIX}/bin/gructl

clean:
	rm -f bin/gructl

.PHONY: get build test integration install uninstall clean
