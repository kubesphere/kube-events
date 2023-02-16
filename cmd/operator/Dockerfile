FROM golang:1.20 as builder
WORKDIR /workspace

ENV GOPROXY=https://goproxy.cn,direct

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd/ cmd/
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GO111MODULE=on go build -a -o ./bin/kube-events-operator ./cmd/operator/main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM kubesphere/distroless-static:nonroot
WORKDIR /
COPY --from=builder /workspace/bin/kube-events-operator /usr/local/bin/kube-events-operator
USER nonroot:nonroot

ENTRYPOINT ["kube-events-operator"]
