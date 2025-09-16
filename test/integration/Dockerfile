FROM golang:1.24.0-alpine as builder

WORKDIR /project

RUN apk add build-base

COPY .. .

RUN go mod tidy
RUN go test -v ./...