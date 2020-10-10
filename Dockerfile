FROM golang:buster

COPY . /app
WORKDIR /app

RUN export GOPATH=/Users/$USER/go
RUN export PATH=$GOPATH/bin:$PATH
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN go mod download
RUN make
