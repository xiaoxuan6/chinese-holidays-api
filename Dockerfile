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
ENV HOLIDAY_PORT=80

COPY --from=build-dev /go/src/app/holiday ./holiday

EXPOSE $HOLIDAY_PORT

CMD ["./holiday"]
