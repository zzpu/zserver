FROM alpine:latest
COPY ./configs /data/zserver/configs
COPY ./zserver /data/zserver/zserver
WORKDIR /data/zserver

RUN apk update \
    && apk upgrade \
    && apk add --no-cache  \
    && apk add ca-certificates \
    && apk add update-ca-certificates 2>/dev/null || true \
    && apk add tzdata \
    && ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

CMD ["/data/zserver/zserver", "-conf", "/data/zserver/configs"]