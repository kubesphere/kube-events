FROM golang:1.13.5 as builder
WORKDIR /workspace
COPY . .
ENV GOPROXY=https://goproxy.cn,direct
RUN CGO_ENABLED=0 GO111MODULE=on go build -a -o ./bin/kube-events-exporter ./cmd/exporter/main.go

FROM alpine:3.9
WORKDIR /
COPY --from=builder /workspace/bin/kube-events-exporter /usr/local/bin/kube-events-exporter

RUN adduser -D -g eventer -u 1002 eventer && \
    chown -R eventer:eventer /usr/local/bin/kube-events-exporter
USER eventer

ENTRYPOINT ["kube-events-exporter"]
