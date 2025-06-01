VERSION ?= 0.0.1
IMG ?= kube-mcp-server:latest

# CONTAINER_TOOL defines the container tool to be used for building images.
CONTAINER_TOOL ?= docker

.PHONY: all
all: build

##@ Development

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test: fmt vet
	go test -v ./... -coverprofile cover.out

.PHONY: lint
lint:
	# Run golangci-lint
	golangci-lint run --timeout 5m

.PHONY: mcp
mcp: build ## Run the Model Context Protocol (MCP) inspector tool
	npx @modelcontextprotocol/inspector@latest bin/kube-mcp-server

##@ Build

.PHONY: build
build: fmt vet
	go build -o bin/kube-mcp-server main.go

.PHONY: run
run: fmt vet
	go run ./main.go

.PHONY: docker-build
docker-build:
	$(CONTAINER_TOOL) build -t ${IMG} . --load

.PHONY: docker-push
docker-push:
	$(CONTAINER_TOOL) push ${IMG}

# PLATFORMS specifies the target platforms for building images, enabling support for multiple
# architectures. (i.e. make docker-buildx IMG=myregistry/image:0.0.1). To use this option you need to:
# - be able to use docker buildx. More info: https://docs.docker.com/build/buildx/
# - have enabled BuildKit. More info: https://docs.docker.com/develop/develop-images/build_enhancements/
# - be able to push the image to your registry (i.e. if you do not set a valid value via IMG=<myregistry/image:<tag>> then the export will fail)
# To adequately provide solutions that are compatible with multiple platforms, you should consider using this option.
PLATFORMS ?= linux/arm64,linux/amd64
.PHONY: docker-buildx-push
docker-buildx-push:
	# copy existing Dockerfile and insert --platform=${BUILDPLATFORM} into Dockerfile.cross, and preserve the original Dockerfile
	sed -e '1 s/\(^FROM\)/FROM --platform=\$$\{BUILDPLATFORM\}/; t' -e ' 1,// s//FROM --platform=\$$\{BUILDPLATFORM\}/' Dockerfile > Dockerfile.cross
	- $(CONTAINER_TOOL) buildx create --name kube-mcp-server-builder
	$(CONTAINER_TOOL) buildx use kube-mcp-server-builder
	- $(CONTAINER_TOOL) buildx build --push --platform=$(PLATFORMS) --tag ${IMG} -f Dockerfile.cross .
	- $(CONTAINER_TOOL) buildx rm kube-mcp-server-builder
	rm Dockerfile.cross

.PHONY: goreleaser-release-snapshot
goreleaser-release-snapshot: ## Build and run release in snapshot mode
	goreleaser release --snapshot --clean
