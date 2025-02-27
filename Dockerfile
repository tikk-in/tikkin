FROM golang:1.24-alpine as builder
WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o tikkin

FROM alpine:latest

RUN apk update && apk upgrade && apk add --no-cache ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=builder /app/tikkin /app/tikkin
COPY --from=builder /app/template/ /app/template/
COPY --from=builder /app/docs/swagger.json /app/docs/swagger.json
COPY --from=builder /app/docs/swagger.yaml /app/docs/swagger.yaml

ENTRYPOINT ["/app/tikkin"]
