FROM golang:alpine

VOLUME ["/dep"]
ENTRYPOINT ["/go/bin/diligent"]

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /go/src/github.com/senseyeio/diligent
COPY . .

RUN go install github.com/senseyeio/diligent/cmd/diligent

WORKDIR /dep

