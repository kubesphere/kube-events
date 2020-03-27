FROM golang:1.13.5 as builder
WORKDIR /workspace
COPY . .
RUN CGO_ENABLED=0 GO111MODULE=on go build -a -o ./bin/kube-events-ruler ./cmd/ruler/main.go

FROM alpine
WORKDIR /
COPY --from=builder /workspace/bin/kube-events-ruler /usr/local/bin/kube-events-ruler
ENTRYPOINT ["kube-events-ruler"]
