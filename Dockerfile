# base go image
FROM golang:1.22-alpine as builder

RUN mkdir /app

COPY ./ /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o captibidadapter ./

RUN chmod +x /app/captibidadapter

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/captibidadapter /app
COPY --from=builder /app/static /app/static
COPY --from=builder /app/config.yaml /app

WORKDIR /app

CMD [ "./captibidadapter" ]