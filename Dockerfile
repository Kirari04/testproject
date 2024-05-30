FROM golang:alpine3.20 as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./main.go

FROM haproxy:2.3
WORKDIR /app

COPY --from=builder /app/main /app/main
COPY haproxy/haproxy.cfg /app/haproxy/haproxy.cfg

ENTRYPOINT ["/app/main", "serve"]