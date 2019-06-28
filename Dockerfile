FROM golang:alpine as builder
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0
WORKDIR /go/src/github.com/hassieswift621/discord-hackweek-2019
ADD . .
RUN apk add --no-cache git
RUN go get ./...
RUN go build -a -installsuffix cgo -o bot.out ./main

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/hassieswift621/discord-hackweek-2019/bot.out /bot
CMD ./bot