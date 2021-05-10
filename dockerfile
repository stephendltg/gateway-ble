FROM golang:1.16.4-alpine AS builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/gateway-ble /app/
WORKDIR /app
CMD ["./gateway-ble"]