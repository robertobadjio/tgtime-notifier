FROM golang:1.23-alpine AS build

RUN apk update && apk add tzdata

WORKDIR /app

COPY ../../go.mod go.sum ./
RUN go mod download

COPY ../../cmd/notifier/notifier.go ./
COPY ../../internal internal
COPY ../../cmd cmd

RUN CGO_ENABLED=0 GOOS=linux go build -o /tgtime-notifier

FROM scratch

COPY --from=build /usr/share/zoneinfo/Europe/Moscow /usr/share/zoneinfo/Europe/Moscow
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /tgtime-notifier .
COPY --from=ghcr.io/tarampampam/curl:8.6.0 /bin/curl /bin/curl

ENV TZ Europe/Moscow

CMD ["/tgtime-notifier", "--config", "config"]