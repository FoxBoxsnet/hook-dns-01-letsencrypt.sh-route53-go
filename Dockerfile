FROM golang:1.6.3-alpine

MAINTAINER FoxBoxsnet

RUN apk add --no-cache git

COPY build-golang.sh /go

WORKDIR /go
ENTRYPOINT [ "/build-golang.sh" ]
CMD ["exit"]