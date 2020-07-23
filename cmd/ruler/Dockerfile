FROM golang:1.13.5 as builder
WORKDIR /workspace
COPY . .
RUN CGO_ENABLED=0 GO111MODULE=on go build -a -o ./bin/kube-events-ruler ./cmd/ruler/main.go

FROM alpine:3.9
WORKDIR /
COPY --from=builder /workspace/bin/kube-events-ruler /usr/local/bin/kube-events-ruler

RUN adduser -D -g eventer -u 1002 eventer && \
    chown -R eventer:eventer /usr/local/bin/kube-events-ruler
USER eventer

ENTRYPOINT ["kube-events-ruler"]
