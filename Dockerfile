FROM gcr.io/distroless/static:latest

LABEL org.opencontainers.image.title=Mercury
LABEL org.opencontainers.image.description="A Planet-style feed aggregator"
LABEL org.opencontainers.image.vendor="Keith Gaughan"
LABEL org.opencontainers.image.licenses=MIT
LABEL org.opencontainers.image.url=https://github.com/kgaughan/mercury
LABEL org.opencontainers.image.source=https://github.com/kgaughan/mercury
LABEL org.opencontainers.image.documentation=https://kgaughan.github.io/mercury/

COPY mercury /

VOLUME ["/data", "/config"]

ENTRYPOINT ["/mercury"]
CMD ["--config", "/config/mercury.toml"]
