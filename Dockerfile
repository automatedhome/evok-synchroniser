FROM arm32v7/golang:stretch

COPY qemu-arm-static /usr/bin/
WORKDIR /go/src/github.com/automatedhome/evok-synchroniser
COPY . .
RUN go build -o synchroniser cmd/main.go

FROM arm32v7/busybox:1.30-glibc

COPY --from=0 /go/src/github.com/automatedhome/evok-synchroniser/synchroniser /usr/bin/synchroniser

ENTRYPOINT [ "/usr/bin/synchroniser" ]
