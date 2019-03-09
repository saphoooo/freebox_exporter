FROM golang:1.7.3 as builder
WORKDIR /
RUN go get -d -v golang.org/x/net/html  
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM scratch  
COPY --from=builder app /
CMD ["/app"]
