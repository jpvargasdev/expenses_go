FROM golang:1.23-alpine as builder

LABEL maintainer="Juan Vargas <vargasm.jp@gmail.com>"

WORKDIR /app

ENV CGO_ENABLED=1
ENV GOOS=linux

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o guilliman cmd/server/main.go

FROM alpine:3.18

WORKDIR /app

RUN apk add --no-cache ca-certificates musl

COPY --from=builder /app/guilliman .

COPY --from=builder /app/.env ./.env

EXPOSE 8080

CMD ["./guilliman"]
