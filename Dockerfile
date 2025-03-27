FROM golang:alpine AS builder

COPY . /src/doc-render
WORKDIR /src/doc-render

RUN set -ex \
 && apk add --update \
      git \
      make \
      nodejs \
      yarn \
 && make frontend_prod \
 && go install \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly \
      -modcacherw \
      -trimpath


FROM alpine:latest

LABEL maintainer="Knut Ahlers <knut@ahlers.me>"

RUN set -ex \
 && apk --no-cache add \
      ca-certificates

COPY --from=builder /go/bin/doc-render /usr/local/bin/doc-render

EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/doc-render"]
CMD ["--"]

# vim: set ft=Dockerfile:
