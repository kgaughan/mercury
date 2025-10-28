FROM gcr.io/distroless/static:latest

LABEL org.opencontainers.image.title=Mercury \
    org.opencontainers.image.description="A Planet-style feed aggregator" \
    org.opencontainers.image.vendor="Keith Gaughan" \
    org.opencontainers.image.licenses=MIT \
    org.opencontainers.image.url=https://github.com/kgaughan/mercury \
    org.opencontainers.image.source=https://github.com/kgaughan/mercury \
    org.opencontainers.image.documentation=https://kgaughan.github.io/mercury/

COPY mercury /

VOLUME ["/data", "/config"]

ENTRYPOINT ["/mercury"]
CMD ["--config", "/config/mercury.toml"]
