# TODOS:
#    - perform proto builds in this file
FROM golang:1.17-alpine as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
# COPY Makefile ./
COPY pkg/ ./pkg
COPY cmd/ ./cmd

# RUN apk --no-cache add make
# RUN make server

RUN go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
CMD ["./server", "-http-port", "8080", "-grpc-port", "9090"]