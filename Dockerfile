FROM golang:1.9.2 AS build-env
LABEL maintainer "szyn <aqr.aqua@gmail.com>"

WORKDIR /go/src/github.com/szyn/mog
COPY . .

RUN go get -u github.com/golang/dep/cmd/dep && \
    dep ensure && \
    CGO_ENABLED=0 go build

FROM busybox:1.27.2
COPY --from=build-env /go/src/github.com/szyn/mog/mog /usr/local/bin/mog
ENTRYPOINT ["/usr/local/bin/mog"]