FROM golang:1.13

WORKDIR /

COPY . .

RUN set -x && \
    go get -d -v . && \
    CGO_ENABLED=0 GOOS=linux GOARM=7 GOARCH=arm go build -a -installsuffix cgo -o app .

FROM scratch

LABEL maintainer="stephane.beuret@gmail.com"

COPY --from=0 app /

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app"]