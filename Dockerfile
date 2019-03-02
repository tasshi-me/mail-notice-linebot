FROM golang:alpine as builder

RUN apk add -U --no-cache \
  git \
  make \
  && rm -rf /var/cache/apk/*

RUN mkdir  %USERPROFILE%\go\src\postgresqlgo
RUN cd %USERPROFILE%\go\src\postgresqlgo
RUN set GOPATH=%USERPROFILE%\go
RUN go get github.com/line/line-bot-sdk-go/linebot

WORKDIR /go
ADD . /go
RUN ["make"]

FROM alpine:latest

WORKDIR /go
COPY --from=builder /go/webapp_linux_amd64 /go/webapp_linux_amd64
CMD ["./webapp_linux_amd64"]