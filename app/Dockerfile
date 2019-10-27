FROM golang:1.13.1 AS build-env

ADD . $GOPATH/src/github.com/tommy-sho/gocon-k8s-server
WORKDIR $GOPATH/src/github.com/tommy-sho/gocon-k8s-server

RUN CGO_ENABLED=0 go build -o app ./main.go

FROM alpine:3.9
COPY --from=build-env /go/src/github.com/tommy-sho/gocon-k8s-server/app /usr/local/bin/app
RUN apk add --no-cache tzdata ca-certificates

ENTRYPOINT ["app"]