FROM golang:1.22-alpine AS builder

WORKDIR /go/src/github.com/IlyaZayats/server
COPY ../.. .

RUN go build -o ./bin/server ./cmd/server

FROM alpine:latest AS runner

COPY --from=builder /go/src/github.com/IlyaZayats/server/bin/server /app/server

RUN apk -U --no-cache add bash ca-certificates \
    && chmod +x /app/server

WORKDIR /app
ENTRYPOINT ["/app/server"]
