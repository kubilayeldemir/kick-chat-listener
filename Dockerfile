FROM golang:1.22.0-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/main .

RUN apk add --no-cache sqlite

CMD ["./main"]