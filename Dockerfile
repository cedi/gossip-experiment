# Build the urlshortener binary
FROM golang:1.22 as builder

WORKDIR /workspace

# Copy the Go Modules manifests
COPY Makefile Makefile
COPY go.mod go.mod
COPY go.sum go.sum

RUN make tools

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download


# Copy the go source
COPY main.go main.go
COPY pkg/ pkg/
COPY cmd/ cmd/

# Build
RUN make build

# Use distroless as minimal base image to package the urlshortener binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
# TODO: For production re-enable distroless!
#FROM gcr.io/distroless/static:nonroot
FROM alpine:latest

WORKDIR /

COPY --from=builder /workspace/gossip .
