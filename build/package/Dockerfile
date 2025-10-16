FROM golang:alpine AS builder

WORKDIR /user-subscriptions

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY ./. ./

RUN go build -v ./cmd/app

FROM alpine:3

COPY --from=builder /user-subscriptions/app /user-subscriptions/app
COPY --from=builder /user-subscriptions/config /config

EXPOSE 8080

ENTRYPOINT ["/user-subscriptions/app"]