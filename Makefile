# Stolen from https://github.com/segmentio/go-release/blob/master/Makefile
build:
	@mkdir -p build
	@go build -o build/release

release: build
	go-release alexanderGugel nerva --assets build/release

.PHONY: build
