FROM golang:1.8.1

RUN go get github.com/Masterminds/glide \
 && go install github.com/Masterminds/glide \
 && go get -u github.com/mitchellh/gox \
 && go get -u github.com/tcnksm/ghr

WORKDIR /go/src/github.com/szyn/mog

CMD ./_tool/build.sh
