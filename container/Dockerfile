FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod . 
COPY go.sum . 
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY config ./config

RUN go build -o bin/main ./cmd/app/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/main ./bin/main
COPY --from=builder /app/config ./config

# Run the app
CMD ["/app/bin/main"]