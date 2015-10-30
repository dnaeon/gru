PREFIX := /usr/local

build:
	glide install
	go build -o bin/gructl gructl/main.go

test:
	go test

install: build
	install -m 0755 bin/gructl ${PREFIX}/bin/gructl

uninstall:
	rm -f ${PREFIX}/bin/gructl

clean:
	rm -f bin/gructl

.PHONY: build test install uninstall clean
