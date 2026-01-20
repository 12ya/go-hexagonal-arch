FROM golang:1.25.6-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/port-service/main.go

FROM alpine:latest

RUN apk update && apk upgrade

RUN rm -rf /var/cache/apk/* && \
    rm -rf /tmp/*

RUN adduser -D appuser
USER appuser

WORKDIR /app

COPY --from=builder /app/app .

ENV HTTP_ADDR=:8080

EXPOSE 8080

CMD [ "./app" ]