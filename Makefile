# Stolen from https://github.com/segmentio/go-release/blob/master/Makefile
build:
	@mkdir -p build
	@go build -o build/nerva

release: build
	go-release alexanderGugel nerva build/nerva

.PHONY: build
