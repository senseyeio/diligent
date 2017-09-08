FROM golang:alpine

VOLUME ["/dep"]
ENTRYPOINT ["/go/bin/dil"]

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /go/src/github.com/senseyeio/diligent
COPY . .

RUN go get ./...
RUN go install github.com/senseyeio/diligent/cmd/dil

WORKDIR /dep

