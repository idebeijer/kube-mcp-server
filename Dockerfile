FROM golang:1.24-alpine AS builder
ENV CGO_ENABLED=0

ARG TARGETOS=linux
ARG TARGETARCH=amd64

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-s -w" -o /app/kube-mcp-server .

FROM alpine:3.22 AS kubectl-installer

ARG TARGETARCH=amd64

RUN apk add --no-cache ca-certificates curl \
    && KUBECTL_VERSION=$(curl -L -s https://dl.k8s.io/release/stable.txt) \
    && curl -LO "https://dl.k8s.io/release/${KUBECTL_VERSION}/bin/linux/${TARGETARCH}/kubectl" \
    && chmod +x kubectl \
    && mv kubectl /usr/local/bin/kubectl

FROM gcr.io/distroless/static-debian12

WORKDIR /app

COPY --from=kubectl-installer /usr/local/bin/kubectl /usr/local/bin/kubectl

COPY --from=builder /app/kube-mcp-server /app

USER nonroot:nonroot
ENTRYPOINT ["/app/kube-mcp-server"]