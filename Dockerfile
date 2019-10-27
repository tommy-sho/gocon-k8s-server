FROM golang:1.13.1 AS build-env
ADD . $GOPATH/src/github.com/tommy-sho/gocon-k8s-server
WORKDIR $GOPATH/src/github.com/tommy-sho/gocon-k8s-server/
RUN CGO_ENABLED=0 go build -ldflags "-X main.version=$(git rev-parse --verify HEAD)" \
    -o vega-golang-layout ./cmd/app/server.go

FROM alpine:3.10.2
COPY --from=build-env /go/src/github.com/tommy-sho/gocon-k8s-server /usr/local/bin/gocon-k8s-server
RUN apk add --no-cache tzdata ca-certificates

CMD ["gocon-k8s-server"]