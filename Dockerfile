FROM golang:1.6.3-alpine

MAINTAINER FoxBoxsnet

COPY build-golang.sh /
ENTRYPOINT [ "/build-golang.sh" ]

