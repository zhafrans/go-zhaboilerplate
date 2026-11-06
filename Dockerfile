# Builder
FROM golang:1.21-alpine as builder

RUN apk update && apk upgrade && \
    apk --update add git make bash build-base

WORKDIR /app

COPY . .

RUN make build

# Distribution
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata ca-certificates && \
    mkdir /app 

WORKDIR /app 

EXPOSE 9090

COPY --from=builder /app/engine /app/

CMD ["/app/engine"]
