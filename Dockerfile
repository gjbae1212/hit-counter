FROM golang:1.14.5-alpine AS builder

WORKDIR /go/src/github.com/gjbae1212/hit-counter

RUN go env -w GO111MODULE="on"

# copy go.mod go.sum
COPY ./go.mod ./go.sum ./

# download Library
RUN go mod download

# copy all
COPY ./ ./

RUN CGO_ENABLED=0 go build -a -ldflags "-w -s" -o /go/bin/hit-counter

# Minimize a docker image
FROM gcr.io/distroless/base:latest

COPY --from=builder /go/bin/hit-counter /go/bin/hit-counter

COPY --from=builder /go/src/github.com/gjbae1212/hit-counter/public /public

COPY --from=builder /go/src/github.com/gjbae1212/hit-counter/view /go/src/github.com/gjbae1212/hit-counter/view

CMD ["/go/bin/hit-counter"]
