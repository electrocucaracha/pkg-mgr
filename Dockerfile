FROM golang:1.13-buster as builder

WORKDIR /go/src/github.com/electrocucaracha/pkg-mgr
COPY . .

ENV GO111MODULE "on"
ENV CGO_ENABLED "0"
ENV GOOS "linux"
ENV GOARCH "amd64"

RUN go build -v -tags netgo -installsuffix netgo -o /bin/pkg_mgr cmd/server/main.go

FROM gcr.io/distroless/base:nonroot
MAINTAINER Victor Morales <electrocucaracha@gmail.com>

ENV PKG_DEBUG "false"
ENV PKG_SQL_ENGINE "sqlite"
ENV PKG_DB_USERNAME ""
ENV PKG_DB_PASSWORD ""
ENV PKG_DB_HOSTNAME ""
ENV PKG_DB_DATABASE "pkg_db"
ENV PKG_SCRIPTS_PATH "./scripts"
ENV PKG_MAIN_FILE "./install.sh"

LABEL io.k8s.display-name="cURL Package Manager"
EXPOSE 3000

COPY --from=builder /bin/pkg_mgr /pkg_mgr

ENTRYPOINT ["/pkg_mgr"]
CMD ["serve"]
