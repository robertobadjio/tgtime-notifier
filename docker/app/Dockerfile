FROM golang:1.22-alpine as build

RUN apk update && apk add tzdata

WORKDIR /app

COPY ../../go.mod go.sum ./
RUN go mod download

RUN ls -la

COPY ../../cmd/notifier/notifier.go ./
COPY ../../internal internal
COPY ../../cmd cmd
COPY ../../.env /

RUN CGO_ENABLED=0 GOOS=linux go build -o /tgtime-notifier

FROM scratch

COPY --from=build /usr/share/zoneinfo/Europe/Moscow /usr/share/zoneinfo/Europe/Moscow
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /tgtime-notifier .
COPY --from=build /.env .

ENV TZ Europe/Moscow

EXPOSE 8080

CMD ["/tgtime-notifier", "--config", "config"]