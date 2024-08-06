FROM golang:1.22.5-alpine3.20 AS build-dev

WORKDIR /go/src/app

COPY . .

RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod download && \
    apk add --no-cache upx && \
    go build -trimpath  -ldflags="-s -w" -o holiday . && \
    [ -e /usr/bin/upx ] && upx holiday || echo

FROM alpine

RUN apk update && \
    apk add --no-cache tzdata

ENV TZ=Asia/Shanghai

COPY --from=build-dev /go/src/app/holiday ./holiday

EXPOSE 80

CMD ["./holiday"]