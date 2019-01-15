FROM golang:alpine as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GO111MODULE=on

RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

RUN adduser -D -g '' appuser

ADD . ${GOPATH}/src/app/
WORKDIR ${GOPATH}/src/app

RUN go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/speedtest_exporter

FROM scratch
ARG VCS_REF
ARG BUILD_DATE

LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.vcs-url="https://github.com/nlamirault/speedtest_exporter" \
      org.label-schema.docker.cmd="docker run tnwhitwell/speedtest_exporter" \
      org.label-schema.schema-version="1.0" \
      maintainer="tom@whi.tw"

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /go/bin/speedtest_exporter /go/bin/speedtest_exporter

USER appuser

ENTRYPOINT [ "/go/bin/speedtest_exporter" ]
