FROM golang:1.11.0-alpine3.8 AS src

VOLUME ["/test-results"]
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh gcc libc-dev

WORKDIR /go/src/github.com/senseyeio/diligent
COPY . .

FROM src as run

VOLUME ["/dep"]
WORKDIR /dep
ENTRYPOINT ["/go/bin/diligent"]
RUN go install github.com/senseyeio/diligent/cmd/diligent
