FROM golang:alpine3.16
WORKDIR app

COPY main.go /app/main.go
RUN go mod init main
RUN go get -t -v github.com/rabbitmq/amqp091-go
RUN go build -o /app/gobuild/RabbitMSQ /app/main.go
ENTRYPOINT ["/app/gobuild/RabbitMSQ", "HttpServerORIG"]
