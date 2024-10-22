# Use the official Golang image to build the application
FROM golang:1.23.2-alpine3.20 AS builder

# Set the working directory
WORKDIR /app

COPY . .

RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz

# Run Stage
FROM alpine:3.20
WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY config.example.yaml ./config.yaml
COPY start.sh .
COPY wait-for.sh .
COPY db/migrations ./migrations

# Make the scripts executable
RUN chmod +x /app/start.sh /app/wait-for.sh

EXPOSE 8080

CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]
