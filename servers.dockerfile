FROM golang:alpine3.16 As base
WORKDIR app
EXPOSE 3333
COPY . .
RUN chmod 777 .
RUN go build -o gobuild/devopstuniapp cmd/main/main.go
ENV mode = ""

#ENTRYPOINT ["gobuild/devopstuniapp", "HttpServ"]
ENTRYPOINT gobuild/devopstuniapp $mode
