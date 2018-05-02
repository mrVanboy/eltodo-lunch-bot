FROM golang:1.10 AS build-env
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /bin/dep
ADD ./src /go/src
RUN chmod +x /bin/dep && \
    cd /go/src/eltodo-lunch-bot && \
    dep ensure -vendor-only && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o goapp

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/eltodo-lunch-bot/goapp /app
RUN apk --no-cache add poppler-utils tzdata  ca-certificates
ENTRYPOINT ./goapp
