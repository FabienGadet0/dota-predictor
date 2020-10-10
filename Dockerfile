FROM golang:buster

COPY . /app
WORKDIR /app

RUN go mod install
RUN make