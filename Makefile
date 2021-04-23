NAME    := $(shell basename $(CURDIR))
VERSION := 0.0.1

ifeq (,$(wildcard ../../.git/HEAD))
REVISION := ${GIT_SHA1_HASH}
else
REVISION := $(shell git rev-parse --short HEAD)
endif

SRCS := $(shell find $(CURDIR) -type f -name '*.go')

GOOS   := linux
GOARCH := amd64

LDFLAGS_NAME     := -X "main.name=$(NAME)"
LDFLAGS_VERSION  := -X "main.version=v$(VERSION)"
LDFLAGS_REVISION := -X "main.revision=$(REVISION)"
LDFLAGS          := -ldflags '-s -w $(LDFLAGS_NAME) $(LDFLAGS_VERSION) $(LDFLAGS_REVISION) -extldflags -static'

.PHONY: all
all: $(NAME)

.PHONY: $(NAME)
$(NAME): $(CURDIR)/bin/$(NAME)
$(CURDIR)/bin/$(NAME): $(SRCS)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $@

$(CURDIR)/bin/$(NAME).zip: $(CURDIR)/bin/$(NAME)
	cd $(CURDIR)/bin && zip $@ $(NAME)

.PHONY: run
run: $(CURDIR)/bin/$(NAME)
	$(CURDIR)/bin/$(NAME) -token ${GITHUB_TOKEN} -owner ${GITHUB_OWNER} -repos ${GITHUB_REPOS} -workflow ${GITHUB_WORKFLOW}

.PHONY: test
test:
	go test -v

.PHONY: clean
clean:
	rm -rf $(CURDIR)/bin
