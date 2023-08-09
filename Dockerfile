FROM golang:1.21 AS builder

COPY . /builder
WORKDIR /builder

RUN make build

FROM alpine:3

WORKDIR /app

COPY --from=builder /builder .

ENTRYPOINT ["/app/vanity"]
