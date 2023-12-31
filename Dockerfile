FROM golang:1.21-alpine

ARG DEBIAN_FRONTEND=noninteractive

WORKDIR /app/
ADD go.mod /app/
ADD go.sum /app
ADD Makefile /app/
ADD ./src/* /app/src/
ADD ./data/* /app/data/

RUN go env -w GO111MODULE=on
RUN go env -w  GOPROXY=https://goproxy.io,direct
RUN go build -o bot ./src

CMD WORDLIST_PATH="./data/wordlist.txt" ./bot