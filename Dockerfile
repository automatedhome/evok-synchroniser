FROM arm32v7/golang:stretch as builder

COPY qemu-arm-static /usr/bin/
WORKDIR /go/src/github.com/automatedhome/evok-synchroniser
COPY . .
RUN make build

FROM arm32v7/busybox:1.30-glibc

COPY --from=builder /go/src/github.com/automatedhome/evok-synchroniser/evok-synchroniser /usr/bin/evok-synchroniser
COPY --from=builder /go/src/github.com/automatedhome/evok-synchroniser/config.yaml /config.yaml

ENTRYPOINT [ "/usr/bin/evok-synchroniser" ]
