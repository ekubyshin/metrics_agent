current_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))

server:
	go build -o cmd/server/server cmd/server/*.go
	chmod +x cmd/server/server

agent:
	go build -o cmd/agent/agent cmd/agent/*.go
	chmod +x cmd/agent/agent

.PHONY: ytest
ytest:
	cd cmd/agent && \
	metricstest metricstest -test.v -test.run=^TestIteration2[AB]*$ \
            -source-path=. \
            -agent-binary-path=cmd/agent/agent

.PHONY: test
test:
	go test ./...