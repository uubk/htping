FROM golang:1.13-buster AS builder

WORKDIR /go/src/app
COPY go.* /go/src/app/
# Try to speed up subsequent builds
RUN go mod download
COPY cmd cmd
COPY pkg pkg
RUN go build -ldflags "-linkmode external -extldflags -static" -a github.com/uubk/htping/cmd/htping

FROM scratch
COPY --from=builder /go/src/app/htping /htping
COPY static/dist/ /static
ENTRYPOINT ["/htping"]