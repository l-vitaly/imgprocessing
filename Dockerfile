FROM lvitaly/golang-upx:latest AS build_base

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

FROM build_base AS builder

WORKDIR /app

COPY . .

RUN	make build; \
    strip --strip-unneeded service; \
    upx service

FROM alpine

ARG BINARY

RUN apk --update add ca-certificates

COPY ./${BINARY} /bin/service

EXPOSE 9000

ENTRYPOINT ["/bin/service"]