current_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))

LOCAL_BIN=$(CURDIR)/bin

include bin-deps.mk

server:
	go build -o cmd/server/server cmd/server/*.go
	chmod +x cmd/server/server

agent:
	go build -o cmd/agent/agent cmd/agent/*.go
	chmod +x cmd/agent/agent

run-agent:
	go run cmd/agent/main.go

run-server:
	go run cmd/server/main.go

.PHONY: test
test:
	go test ./... -v -cover 

.PHONY: lint
lint: $(GOLANGCI_BIN) ## go lint
	$(GOLANGCI_BIN) run --fix ./...

.PHONE: ytest
ytest: $(METRICSTEST)
	$(METRICSTEST) '-test.v' '-test.run=^TestIteration3[AB]*$\' '-source-path=.' '-agent-binary-path=cmd/agent/agent' '-binary-path=cmd/server/server'