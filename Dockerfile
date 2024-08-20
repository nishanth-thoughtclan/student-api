FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd/ cmd/

WORKDIR /app/cmd/app

RUN go build -o main .

FROM debian:bullseye-slim

COPY --from=builder /app/cmd/app/main /app/main

EXPOSE 8080

CMD ["/app/main"]