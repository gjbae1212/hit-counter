FROM gcr.io/gcp-runtimes/go1-builder:1.12 as builder

ENV GO111MODULE off

ENV GOPATH /go

WORKDIR /go/src/github.com/gjbae1212/hit-counter

RUN mkdir -p /go/src/github.com/gjbae1212/hit-counter

COPY . /go/src/github.com/gjbae1212/hit-counter

RUN /usr/local/go/bin/go build -o /go/bin/hit-counter -tags appenginevm

# Minimize a docker image

FROM gcr.io/distroless/base:latest

COPY --from=builder /go/bin/hit-counter /go/bin/hit-counter

CMD ["/go/bin/hit-counter"]