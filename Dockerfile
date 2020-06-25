FROM golang:1.13-alpine AS builder
WORKDIR /go/src/bowot
COPY . .
RUN go install -v ./cmd/bowot

FROM alpine
RUN apk add --no-cache ca-certificates ffmpeg wget python2
RUN wget https://yt-dl.org/downloads/latest/youtube-dl -O /usr/local/bin/youtube-dl && chmod a+rx /usr/local/bin/youtube-dl	
WORKDIR /root/
COPY --from=builder /go/bin/bowot .
COPY config.yaml .
CMD ["./bowot"]