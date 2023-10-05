current_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))

LOCAL_BIN=$(CURDIR)/bin

include bin-deps.mk

server:
	go build -o cmd/server/server cmd/server/*.go
	chmod +x cmd/server/server

agent:
	go build -o cmd/agent/agent cmd/agent/*.go
	chmod +x cmd/agent/agent

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint: $(GOLANGCI_BIN) ## go lint
	$(GOLANGCI_BIN) run --fix ./...