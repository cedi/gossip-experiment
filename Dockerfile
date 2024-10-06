# Build the urlshortener binary
FROM golang:1.23 AS builder

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY pkg/ pkg/
COPY cmd/ cmd/

# Build
RUN go build -o build/gossip ./main.go

# Use distroless as minimal base image to package the urlshortener binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
# TODO: For production re-enable distroless!
# FROM gcr.io/distroless/static:nonroot AS final
# FROM alpine:latest

WORKDIR /

# COPY --from=builder /workspace/build/gossip .
EXPOSE 7946

# ENTRYPOINT [ "/gossip", "memberlist", "join", "--config", "lan", "-port", "7946", "gossip.gossip.svc.cluster.local:7946"]