# syntax=docker/dockerfile:1

FROM golang:1.25.5-alpine

RUN apk add --no-cache git
RUN apk add build-base

WORKDIR /app

COPY ./server/go.mod .
COPY ./server/go.sum .


RUN go mod download

COPY ./client/build ./client
COPY ./server .

RUN go build .

EXPOSE 8080

CMD [ "./tinylytics" ]