FROM golang:alpine AS src

VOLUME ["/dep"]

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /go/src/github.com/senseyeio/diligent
COPY . .

FROM src as run

WORKDIR /dep
ENTRYPOINT ["/go/bin/diligent"]
RUN go install github.com/senseyeio/diligent/cmd/diligent
