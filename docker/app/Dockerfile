FROM golang:1.22-alpine as build

RUN apk update && apk add tzdata

WORKDIR /app

COPY ../../go.mod go.sum ./
RUN go mod download

COPY ../../cmd/notifier ./cmd/notifier
COPY ../../.env /

RUN CGO_ENABLED=0 GOOS=linux go build -o /notifier -v ./cmd/officetime/notifier.go

FROM scratch

COPY --from=build /usr/share/zoneinfo/Europe/Moscow /usr/share/zoneinfo/Europe/Moscow
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /notifier .
COPY --from=build /.env .

ENV TZ Europe/Moscow

EXPOSE 8080

CMD ["/notifier", "--config", "config"]