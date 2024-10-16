BINDIR=../../bin
APPNAME ?= api
DEST = $(BINDIR)/$(APPNAME)
GO=GOOS=linux CGO_ENABLED=0

.PHONY: build
build:
	$(GO) go build -ldflags "-X main.version=$$VERSION" -o $(DEST) *.go

.PHONY: run
run: build
	@$(DEST)