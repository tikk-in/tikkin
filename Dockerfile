FROM golang:1.22-alpine as builder
WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o tikkin

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/tikkin /app/tikkin

ENTRYPOINT ["/app/tikkin"]
