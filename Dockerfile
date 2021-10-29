FROM golang:1.17-alpine AS src

ENV GO111MODULE=on
VOLUME ["/test-results"]
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh gcc libc-dev

WORKDIR /go/src/github.com/senseyeio/diligent
COPY . .

FROM src as run

VOLUME ["/dep"]
RUN go install -mod vendor github.com/senseyeio/diligent/cmd/diligent
WORKDIR /dep
ENTRYPOINT ["/go/bin/diligent"]
