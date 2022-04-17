FROM golang:alpine

LABEL maintainer="nrmadi02 <nrmadi02@gmail.com>"

RUN mkdir /app

ADD . /app/

WORKDIR /app

RUN go get -d

RUN go build -o main .

CMD ["/app/main"]

EXPOSE 8080