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
	go clean -testcache && go test ./... -v -cover 

.PHONY: lint
lint: $(GOLANGCI_BIN) ## go lint
	$(GOLANGCI_BIN) run --fix ./...
## '-test.run=^TestIteration11$\' \

DB_DSN:=host=localhost user=postgres password=password dbname=metrics_agent port=5432 sslmode=disable
.PHONY: ytest
ytest: $(METRICSTEST) server agent
	$(METRICSTEST) '-test.v' \
	'-source-path=.' \
	'-agent-binary-path=cmd/agent/agent' \
	'-binary-path=cmd/server/server' \
	'-server-port=8080' \
	'-file-storage-path=internal/storage/test/test2.json' \
	'-database-dsn=$(DB_DSN)'

.PHONY: goose-create
goose-create: $(GOOSE)
	env GOOSE_MIGRATION_DIR=./internal/storage/migrations $(GOOSE) create init sql

.PHONY: goose-up
goose-up: $(GOOSE)
	env GOOSE_MIGRATION_DIR=./internal/storage/migrations $(GOOSE) postgres "$(DB_DSN)" up

.PHONY: goose-down
goose-down: $(GOOSE)
	env GOOSE_MIGRATION_DIR=./internal/storage/migrations $(GOOSE) postgres "$(DB_DSN)" down

.PHONY: goose-validate
goose-validate: $(GOOSE)
	env GOOSE_MIGRATION_DIR=./internal/storage/migrations $(GOOSE) validate