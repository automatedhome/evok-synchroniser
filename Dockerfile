FROM arm32v7/golang:stretch

COPY qemu-arm-static /usr/bin/
WORKDIR /go/src/github.com/automatedhome/evok-synchroniser
COPY . .
RUN make build

FROM arm32v7/busybox:1.30-glibc

COPY --from=0 /go/src/github.com/automatedhome/evok-synchroniser/evok-synchroniser /usr/bin/evok-synchroniser

ENTRYPOINT [ "/usr/bin/evok-synchroniser" ]
