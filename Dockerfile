FROM golang:1.26.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o api ./cmd/api
RUN go build -o analytics ./cmd/analytics

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/api ./api
COPY --from=builder /app/analytics ./analytics
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./api"]
