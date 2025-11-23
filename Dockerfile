FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/app

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/server /app/server

RUN chmod +x /app/server

CMD ["/app/server"]
