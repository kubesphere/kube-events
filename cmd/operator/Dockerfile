FROM golang:1.13.5 as builder
WORKDIR /workspace
COPY . .
RUN CGO_ENABLED=0 GO111MODULE=on go build -a -o ./bin/kube-events-operator ./cmd/operator/main.go

FROM alpine
WORKDIR /
COPY --from=builder /workspace/bin/kube-events-operator /usr/local/bin/kube-events-operator
ENTRYPOINT ["kube-events-operator"]
