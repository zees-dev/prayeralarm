## binaryless/native audio output
FROM golang:alpine as builder
LABEL maintainer="github.com/zees-dev"

RUN apk add --update --no-cache git make build-base alsa-lib-dev nodejs npm

WORKDIR /go/src/app

# Make use of docker-layer caching - faster builds
COPY go.* ./ *.go ./ mp3/ ./ 
RUN go mod download
RUN CGO_ENABLED=1 go build -o /go/src/app/app

# Install and build client dependencies
RUN cd client && npm install && npm run build


FROM alpine:latest

RUN apk add --update --no-cache omxplayer tini tzdata

WORKDIR /app

COPY --from=builder /go/src/app/app /app/app
COPY --from=builder /go/src/app/mp3 /app/mp3
COPY --from=builder /go/src/app/client/public /app/client/public

ENTRYPOINT ["/sbin/tini", "--", "/app/app"]
CMD ["@"]
