GOLANGCI_BIN=$(LOCAL_BIN)/golangci-lint
$(GOLANGCI_BIN):
	GOBIN=$(LOCAL_BIN) curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.53.3