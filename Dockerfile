FROM golang:alpine

ADD . /go/src/github.com/Luzifer/gh-NeutronStars
WORKDIR /go/src/github.com/Luzifer/gh-NeutronStars

RUN set -ex \
 && apk add --update git ca-certificates \
 && go install -ldflags "-X main.version=$(git describe --tags || git rev-parse --short HEAD || echo dev)" \
 && apk del --purge git

EXPOSE 3000

ENTRYPOINT ["/go/bin/gh-NeutronStars"]
CMD ["--"]
