PREFIX := /usr/local

build:
	go build -o bin/gructl -v

test:
	go test ./...

install: build
	install -m 0755 bin/gructl ${PREFIX}/bin/gructl

uninstall:
	rm -f ${PREFIX}/bin/gructl

clean:
	rm -f bin/gructl

.PHONY: build test install uninstall clean
