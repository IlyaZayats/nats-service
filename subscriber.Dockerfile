FROM golang:1.22-alpine AS builder

WORKDIR /go/src/github.com/IlyaZayats/subscriber
COPY ../.. .

RUN go build -o ./bin/subscriber ./cmd/subscriber

FROM alpine:latest AS runner

COPY --from=builder /go/src/github.com/IlyaZayats/subscriber/bin/subscriber /app/subscriber

RUN apk -U --no-cache add bash ca-certificates \
    && chmod +x /app/subscriber

WORKDIR /app
ENTRYPOINT ["/app/subscriber"]
