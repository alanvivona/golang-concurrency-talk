FROM golang:latest AS build

COPY ./ /go/src/app
WORKDIR /go/src/app

CMD ./run