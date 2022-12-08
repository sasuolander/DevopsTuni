FROM golang:alpine3.16
WORKDIR app
EXPOSE 3333
COPY . .
RUN chmod 777 .
RUN go build -o gobuild/devopstuniapp cmd/main/main.go
ENTRYPOINT ["gobuild/devopstuniapp", "HttpServerORIG"]
