FROM golang:1.13.5 as builder
WORKDIR /workspace
COPY . .
RUN CGO_ENABLED=0 GO111MODULE=on go build -a -o ./bin/kube-events-exporter ./cmd/exporter/main.go

FROM alpine
WORKDIR /
COPY --from=builder /workspace/bin/kube-events-exporter /usr/local/bin/kube-events-exporter
ENTRYPOINT ["kube-events-exporter"]
