FROM golang:1.22-alpine as builder
WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o tikkin

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/tikkin /app/tikkin
COPY --from=builder /app/docs/swagger.json /app/docs/swagger.json
COPY --from=builder /app/docs/swagger.yaml /app/docs/swagger.yaml

ENTRYPOINT ["/app/tikkin"]
