FROM alpine:latest

LABEL org.opencontainers.image.title=Mercury
LABEL org.opencontainers.image.description="A Planet-style feed aggregator"
LABEL org.opencontainers.image.vendor="Keith Gaughan"
LABEL org.opencontainers.image.licenses=MIT
LABEL org.opencontainers.image.url=https://github.com/kgaughan/mercury
LABEL org.opencontainers.image.source=https://github.com/kgaughan/mercury
LABEL org.opencontainers.image.documentation=https://kgaughan.github.io/mercury/

RUN apk --no-cache add ca-certificates tzdata
COPY mercury .
USER nobody
ENTRYPOINT ["/mercury"]
