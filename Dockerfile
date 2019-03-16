FROM golang:1.11.5

LABEL maintainer="stephane.beuret@gmail.com"

WORKDIR /

RUN go get -d -v golang.org/x/net/html

RUN go get -d -v github.com/prometheus/client_golang/prometheus

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM scratch  

COPY --from=0 app /

ENTRYPOINT ["/app"]
