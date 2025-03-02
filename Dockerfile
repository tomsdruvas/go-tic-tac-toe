FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o app

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 8080 8090

CMD ["./app"]
