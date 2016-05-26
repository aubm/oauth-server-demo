FROM golang:1.6-alpine

RUN apk add --update --no-cache \
    git \
    && rm -rf /var/cache/apk/*

ADD . /go/src/github.com/aubm/oauth-server-demo

RUN go get github.com/aubm/oauth-server-demo/... && \
    go install github.com/aubm/oauth-server-demo

ENTRYPOINT /go/bin/oauth-server-demo

EXPOSE 8080