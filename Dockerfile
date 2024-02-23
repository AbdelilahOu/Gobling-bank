# build the app
FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz


# run the app
FROM alpine
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate.linux-amd64 ./migrate
COPY db/migrations ./migrations
COPY app.env .
COPY start.sh .
COPY wait-for.sh .

# 
EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]