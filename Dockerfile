FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates && update-ca-certificates

COPY . .

RUN go mod tidy

RUN go build -o geoservice main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates && update-ca-certificates

COPY --from=builder /app/geoservice .

COPY --from=builder /app/docs ./docs

EXPOSE 8080

CMD ["./geoservice"]
