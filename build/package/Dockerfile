FROM golang:1.24-alpine AS builder

WORKDIR /user-subscriptions

COPY go.* ./

RUN go mod download && go mod verify

COPY ./. ./

RUN go build -v ./cmd/app

FROM alpine:3

COPY --from=builder /user-subscriptions/app /user-subscriptions/app
COPY --from=builder /user-subscriptions/config /config
COPY --from=builder /user-subscriptions/migration /migration


EXPOSE 8080

ENTRYPOINT ["/user-subscriptions/app"]