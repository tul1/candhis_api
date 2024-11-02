FROM golang:1.23 AS builder

WORKDIR /go/src/app

COPY . .

RUN make download deps-openapi build-api

FROM alpine:latest

WORKDIR /opt/app
COPY --from=builder /go/src/app/bin/api /opt/app/api
COPY --from=builder /go/src/app/conf/api.yml /opt/conf/api.yml

EXPOSE 8080

CMD ["./api", "-config", "/opt/conf/api.yml"]
