.PHONY: build clean test

build:
	go build -o webpanel ./cmd/webpanel

clean:
	rm -f webpanel
	rm -rf dist/

test:
	go test ./...

build-all:
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/webpanel_linux_amd64 ./cmd/webpanel
	GOOS=linux GOARCH=386 go build -o dist/webpanel_linux_i386 ./cmd/webpanel
	GOOS=linux GOARCH=arm64 go build -o dist/webpanel_linux_arm64 ./cmd/webpanel
	GOOS=linux GOARCH=arm GOARM=7 go build -o dist/webpanel_linux_armv7 ./cmd/webpanel
	GOOS=linux GOARCH=arm GOARM=6 go build -o dist/webpanel_linux_armv6 ./cmd/webpanel
