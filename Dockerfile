FROM golang:1.13.5-alpine3.10 AS builder

WORKDIR /build
RUN adduser -u 10001 -D app-runner

ENV GOPROXY https://goproxy.cn
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build  -ldflags "-s -w" -a -o zserver cmd/main.go

FROM alpine:latest AS final
COPY ./configs /data/zserver/configs
COPY --from=builder /build/zserver /data/zserver/zserver
WORKDIR /data/zserver

RUN apk update \
    && apk upgrade \
    && apk add --no-cache  \
    && apk add ca-certificates \
    && apk add update-ca-certificates 2>/dev/null || true \
    && apk add tzdata \
    && ln -fs /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

CMD ["/data/zserver/zserver", "-conf", "/data/zserver/configs"]