FROM golang:alpine

LABEL maintainer="nrmadi02 <nrmadi02@gmail.com>"

RUN apk add git

RUN mkdir /app

RUN go install github.com/swaggo/swag/cmd/swag@latest

RUN swag init -g app/server.go

ADD . /app/

WORKDIR /app

RUN go get -d

RUN go build -o main .

CMD ["/app/main"]

EXPOSE 8080