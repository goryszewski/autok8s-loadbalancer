FROM golang:1.22.0 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o loadbalancer-external .

FROM alpine:3.6

WORKDIR /app

COPY --from=builder /app/loadbalancer-external /app/loadbalancer-external

CMD ["./loadbalancer-external"]
