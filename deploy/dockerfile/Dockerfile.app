# syntax=docker/dockerfile:1

ARG GO_VERSION=1.25

FROM golang:${GO_VERSION}-alpine AS builder

ARG GOPROXY=https://goproxy.cn,direct

WORKDIR /src

RUN apk add --no-cache git ca-certificates

ENV GOPROXY=${GOPROXY} \
    CGO_ENABLED=0 \
    GOOS=linux \
    GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG APP
ARG TYPE
RUN go build -ldflags="-s -w" -o /out/app ./apps/${APP}/${TYPE}

FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata curl netcat-openbsd \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo Asia/Shanghai > /etc/timezone

WORKDIR /app

COPY --from=builder /out/app /app/bin

ARG APP
ARG TYPE
COPY apps/${APP}/${TYPE}/etc/${APP}.yaml /app/conf/config.yaml

ENV TZ=Asia/Shanghai

EXPOSE 8080 8888

ENTRYPOINT ["/app/bin", "-f", "/app/conf/config.yaml"]
