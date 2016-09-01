FROM golang:1.6.3-alpine

MAINTAINER FoxBoxsnet

RUN apk add --no-cache git

COPY build-golang.sh /
ENTRYPOINT [ "/build-golang.sh" ]
CMD ["exit"]