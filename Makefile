##@ Development

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

##@ Build

.PHONY: build
build: fmt vet
	go build -o bin/kube-mcp-server main.go

.PHONY: run
run: fmt vet
	go run ./main.go

.PHONY: mcp
mcp: build
	npx @modelcontextprotocol/inspector@latest bin/kube-mcp-server