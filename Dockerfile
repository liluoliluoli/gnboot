FROM golang:1.23-alpine3.20 AS builder

#ENV GOPROXY=https://goproxy.cn
# CGO_ENABLED=1, need check `ldd --version` is same as builder
ENV CGO_ENABLED=0

RUN apk update && apk add --no-cache git make

COPY . /src
WORKDIR /src
# download first can use docker build cache if go.mod not change
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

COPY .. .
RUN make build

FROM alpine:3.18

RUN apk update && apk add --no-cache bash

COPY --from=builder /src/bin /app

WORKDIR /app

COPY configs/prod /data/conf

CMD ["sh", "-c", "./gnboot -c /data/conf"]
