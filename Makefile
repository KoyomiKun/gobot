HOMEDIR    := $(shell pwd)
OUTDIR     := $(HOMEDIR)/output
GIT_COMMIT := $(shell git rev-parse HEAD)

GO      := $(GO_1_16_BIN)/go
GOROOT  := $(GO_1_16_HOME)
GOPATH  := $(shell $(GO) env GOPATH)
GOMOD   := $(GO) mod
GOBUILD := CGO_ENABLED=0 $(GO) build -ldflags "-s -w -X 'github.com/Koyomikun/gobot/pkg/buildinfo.Commit=$(GIT_COMMIT)'"
GOTEST  := $(GO) test -gcflags="-N -l"
GOPKGS  := $$($(GO) list ./...| grep -vE "vendor")

COVPROF := $(HOMEDIR)/covprof.out  # coverage profile
COVFUNC := $(HOMEDIR)/covfunc.txt  # coverage profile information for each function
COVHTML := $(HOMEDIR)/covhtml.html # HTML representation of coverage profile

all: prepare compile

local:
	CGO_ENABLED=0 go build -o ./gobot -ldflags "-s -w -X 'github.com/Koyomikun/gobot/pkg/buildinfo.Commit=$(GIT_COMMIT)'" cmd/main.go

local-linux:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./gobot -ldflags "-s -w -X 'github.com/Koyomikun/gobot/pkg/buildinfo.Commit=$(GIT_COMMIT)'" cmd/main.go
	which upx &>/dev/null && upx -1 ./gobot

prepare: gomod
	mkdir -p $(OUTDIR)

gomod:
	$(GO) env -w GO111MODULE=on
	$(GO) env -w GONOSUMDB=\*
	$(GOMOD) download

compile:
	$(GOBUILD) -o $(OUTDIR)/gobot cmd/main.go

test: prepare test-case
test-case:
	$(GOTEST) -v -cover $(GOPKGS)

clean:
	rm -rf $(OUTDIR)

# avoid filename conflict and speed up build 
.PHONY: all local prepare compile test clean

