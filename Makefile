NAME      := mog
# for build
VERSION   := v0.1.6
REVISION  := $(shell git rev-parse --short HEAD)
LDFLAGS   := -ldflags="-s -w -X \"main.Version=$(VERSION)\" -X \"main.Revision=$(REVISION)\""
# for dist
DIST_DIRS := find * -type d -exec

.PHONY: setup
## Install dev dependencies
setup:
	go get github.com/golang/dep/cmd/dep
	go get github.com/golang/lint/golint
	go get golang.org/x/tools/cmd/goimports

.PHONY: clean
## Clean resources
clean:
	rm -rf bin/*
	rm -rf dist/*
	rm -rf vendor/*

.PHONY: test
## Run tests
test: deps
	go test -cover -v $(go list ./... | grep -v /vendor/)

.PHONY: install
## Install binary to $GOPATH/bin
install: deps
	go install $(LDFLAGS)

.PHONY: build
## Run build binary to bin
build: deps
	go build $(LDFLAGS) -o bin/$(NAME)

.PHONY: cross-build
## Run cross build
cross-build: deps
	for os in darwin linux windows; do \
		for arch in amd64 386; do \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 \
			go build $(LDFLAGS) -o dist/$$os-$$arch/$(NAME); \
		done; \
	done

.PHONY: dist
## Make dist
dist:
	cd dist && \
	for os in darwin linux windows; do \
		for arch in amd64 386; do \
			cd $$os-$$arch && \
			tar -zcvf ../$(NAME)_$$os\_$$arch.tar.gz . && \
			cd ..; \
		done; \
	done && \
	$(DIST_DIRS) rm -rf {} + && \
	shasum -a 256 * > sha256sums.txt && \
	cd ..

.PHONY: deps
## Install dependencies
deps: setup
	dep ensure

.PHONY: update
## Update all dependencies
update: setup
	dep ensure -update

.PHONY: lint
## Lint
lint: setup
	go vet -v
	for pkg in $$(go list ./... | grep -v /vendor/); do \
		golint --set_exit_status $$pkg || exit $$?; \
	done

.PHONY: fmt
## Format source codes
fmt: setup
	goimports -w
