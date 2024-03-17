FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy && go build cmd/serv.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/serv .
CMD [ "sh", "-c", "./serv"]
