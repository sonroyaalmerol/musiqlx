FROM golang:1.22-alpine AS builder
RUN apk add -U --no-cache \
    build-base \
    ca-certificates \
    git \
    sqlite \
    zlib-dev \
    go

WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN GOOS=linux go build -o musiqlx cmd/musiqlx/musiqlx.go

FROM alpine:3.20
LABEL org.opencontainers.image.source https://github.com/sonroyaalmerol/musiqlx
RUN apk add -U --no-cache \
    ffmpeg \
    mpv \
    ca-certificates \
    tzdata \
    tini \
    shared-mime-info

COPY --from=builder \
    /usr/lib/libgcc_s.so.1 \
    /usr/lib/libstdc++.so.6 \
    /usr/lib/libtag.so.2 \
    /usr/lib/
COPY --from=builder \
    /src/musiqlx \
    /bin/
VOLUME ["/cache", "/data", "/music", "/podcasts"]
EXPOSE 80
ENV TZ ""
ENV MUSIQLX_DB_PATH /data/musiqlx.db
ENV MUSIQLX_LISTEN_ADDR :80
ENV MUSIQLX_MUSIC_PATH /music
ENV MUSIQLX_PODCAST_PATH /podcasts
ENV MUSIQLX_CACHE_PATH /cache
ENV MUSIQLX_PLAYLISTS_PATH /playlists
ENTRYPOINT ["/sbin/tini", "--"]
CMD ["musiqlx"]
