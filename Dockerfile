FROM golang:alpine as builder
LABEL maintainer="github.com/zees-dev"

RUN apk update
RUN apk add --no-cache git make build-base alsa-lib-dev

WORKDIR /go/src/app

# Make use of docker-layer caching - faster builds
COPY go.* ./ *.go ./ mp3/ ./ 
RUN go mod download
RUN CGO_ENABLED=1 go build -o /go/src/app/app


FROM alpine:3.12.3

RUN apk update && apk add --no-cache tzdata pulseaudio tini

WORKDIR /app

COPY --from=builder /go/src/app/app /app/app
COPY --from=builder /go/src/app/mp3 /app/mp3

ENTRYPOINT ["/sbin/tini", "--", "/app/app"]
CMD ["@"]
