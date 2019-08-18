FROM golang:alpine as builder

RUN apk add -U --no-cache \
  git \
  make \
  && rm -rf /var/cache/apk/*

ENV GO111MODULE on
ENV GOFLAGS=-mod=vendor

# Recompile the standard library without CGO
RUN CGO_ENABLED=0 go install -a std

ENV APP_DIR $GOPATH/src/github.com/mshrtsr/mail-notice-linebot/
WORKDIR ${APP_DIR}

COPY . $APP_DIR
RUN ["make"]

FROM alpine:latest
RUN apk add -U --no-cache \
  ca-certificates \
  && update-ca-certificates 2>/dev/null || true \
  && rm -rf /var/cache/apk/*

ENV APP_DIR /go/src/github.com/mshrtsr/mail-notice-linebot/
WORKDIR ${APP_DIR}
COPY --from=builder ${APP_DIR}/webapp_linux_amd64 ${APP_DIR}/webapp_linux_amd64
CMD ["./webapp_linux_amd64"]