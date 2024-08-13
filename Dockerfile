FROM golang:1.22.0 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o loadbalancer-external ./cmd/main.go

FROM ubuntu

WORKDIR /app

COPY --from=builder /app/loadbalancer-external /app/loadbalancer-external
COPY --from=builder /app/haproxy.tpl /app/haproxy.tpl
CMD ["./loadbalancer-external"]
