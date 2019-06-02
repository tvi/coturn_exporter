FROM golang:1.12
RUN go get -u github.com/golang/dep/cmd/dep
WORKDIR /go/src/github.com/tvi/coturn_exporter/
COPY . /go/src/github.com/tvi/coturn_exporter/
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ./cmd/exporter



FROM alpine:3.9
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/tvi/coturn_exporter/exporter .
CMD /root/exporter
