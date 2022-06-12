FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go get
RUN CGO_ENABLED=0 go build -o proxy

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/proxy .
EXPOSE 3000
CMD ["./proxy"]