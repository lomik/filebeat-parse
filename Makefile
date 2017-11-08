GO ?= go
export GOPATH := $(CURDIR)/_vendor

all: build

patch:
	rm -rf data
	rm -rf _vendor/src/github.com/elastic/beats/libbeat/processors/actions/parse
	ln -s ../../../../../../../../parse _vendor/src/github.com/elastic/beats/libbeat/processors/actions/parse
	cp inject_parse.go _vendor/src/github.com/elastic/beats/libbeat/processors/actions/inject_parse.go

clean:
	rm -rf _vendor/src/github.com/elastic/beats/libbeat/processors/actions/parse
	rm -rf _vendor/src/github.com/elastic/beats/libbeat/processors/actions/inject_parse.go
	rm -f filebeat

build: patch
	go build github.com/elastic/beats/filebeat

run: build
	./filebeat -c sample.yaml

test:
	cd parse && go test -bench=.

submodules:
	git submodule sync
	git submodule update --init --recursive
