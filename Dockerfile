FROM golang:1.17 AS builder
WORKDIR /go/src/bot
COPY . .
RUN go mod tidy && \
go build cmd/main.go

FROM ubuntu:20.04 AS bot
WORKDIR /bot
COPY --from=builder /go/src/bot/main .
COPY --from=builder /go/src/bot/.env .

RUN apt update && apt install ca-certificates -y

ENTRYPOINT [ "/bot/main" ]