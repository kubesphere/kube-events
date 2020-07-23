FROM golang:1.13.5 as builder
WORKDIR /workspace
COPY . .
RUN CGO_ENABLED=0 GO111MODULE=on go build -a -o ./bin/kube-events-operator ./cmd/operator/main.go

FROM alpine:3.9
WORKDIR /
COPY --from=builder /workspace/bin/kube-events-operator /usr/local/bin/kube-events-operator

RUN adduser -D -g eventer -u 1002 eventer && \
    chown -R eventer:eventer /usr/local/bin/kube-events-operator
USER eventer

ENTRYPOINT ["kube-events-operator"]
