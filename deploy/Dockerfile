FROM golang:alpine AS builder

RUN mkdir /app
ADD .. /app
WORKDIR /app

RUN go build -o /build

EXPOSE 8088

FROM alpine:latest
WORKDIR /app
COPY --from=builder /build .
RUN mkdir /app/assets
COPY assets /app/assets

CMD ["/app/build"]

