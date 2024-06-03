FROM golang:alpine3.20 as builder

WORKDIR /app
RUN apk update
RUN apk add build-base sqlite
RUN update-ca-certificates

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -ldflags="-w -s" -o main ./main.go

FROM node:20 as vue
WORKDIR /app
RUN npm install -g bun
ENV NODE_ENV=production
COPY vue/ .
RUN echo "VITE_APP_API=" > .env
RUN bun install --frozen-lockfile --production
RUN bun run build

FROM alpine:3.20
WORKDIR /app
RUN apk add sqlite haproxy

COPY --from=builder /app/main /app/main
COPY haproxy/haproxy.cfg /app/haproxy/haproxy.cfg
COPY --from=vue /app/dist /app/dist

ENTRYPOINT ["/app/main", "serve", "--tls"]