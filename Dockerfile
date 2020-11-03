FROM alpine:latest

RUN apk  --no-cache --update upgrade && \
    apk  --no-cache add \
        tini \
        ca-certificates

WORKDIR /opt/

COPY bcache ./bcache

ENTRYPOINT ["/sbin/tini"]
CMD ["/opt/bcache", "-db", "/var/bcache/cache.db", "-rm"]
