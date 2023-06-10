FROM golang:1.18 AS builder

ENV GOPATH /go
ENV APPPATH /repo
COPY . /repo
RUN cd /repo && CGO_ENABLED=0 go build -tags netgo -trimpath -ldflags '-s -w' -o mercury ./cmd/mercury

FROM alpine:latest
COPY --from=builder /repo/mercury /mercury
ENTRYPOINT ["/mercury"]
