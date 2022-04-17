FROM golang:alpine

LABEL maintainer="nrmadi02 <nrmadi02@gmail.com>"

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod ./

RUN go mod tidy

COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o /binary

EXPOSE 8080

CMD ["/binary"]