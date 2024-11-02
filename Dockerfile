FROM golang:1.23-alpine as builder

LABEL maintainer="Juan Vargas <vargasm.jp@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o guilliman cmd/server/main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/guilliman .

COPY --from=builder /app/.env ./.env

EXPOSE 8080

CMD ["./guilliman"]
