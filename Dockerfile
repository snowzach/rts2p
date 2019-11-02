FROM golang:1.13-alpine3.10 as builder

RUN apk add --no-cache git live-media live-media-dev gcc g++ libc-dev && \
    rm -rf /var/cache/apk/*
ENV CGO_ENABLED 1
ENV GOOS linux

WORKDIR /build
ADD . .
RUN go build

FROM alpine:3.10
RUN apk add --no-cache live-media libstdc++ && \
    rm -rf /var/cache/apk/*
WORKDIR /opt/rts2p
COPY --from=builder /build/rts2p /opt/rts2p/rts2p
COPY example.yaml /opt/rts2p/rts2p.yaml
CMD [ "/opt/rts2p/rts2p", "-c", "/opt/rts2p/rts2p.yaml" ]