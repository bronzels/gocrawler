# Build stage
FROM golang:1.13.5-alpine3.11
LABEL "Evan Lin"="evanslin@gmail.com"

RUN apk add --no-cache protobuf ca-certificates git
ENV CGO_ENABLED 0
ENV GOOS linux
ENV TIMEZONE "Asia/Taipei"
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/${TIMEZONE} /etc/localtime && \
    echo $TIMEZONE > /etc/timezone && \
    apk del tzdata
COPY . /go/src/github.com/bronzels/gocrawler/
COPY big.jpg /go/bin/
RUN go install ./src/...

# the final image
FROM alpine:3.11
RUN apk add --no-cache ca-certificates rsync shadow bash

# copy the go binaries from the build image
COPY --from=0 /go/bin /go/bin
WORKDIR /go/bin

