FROM golang:1.17 AS builder
WORKDIR /go/src/bot
COPY . .
RUN go mod tidy && \
go build cmd/main.go

FROM ubuntu:20.04 AS bot
WORKDIR /bot
COPY --from=builder /go/src/bot/main .
COPY --from=builder /go/src/bot/.env .

ENTRYPOINT [ "/bot/main" ]