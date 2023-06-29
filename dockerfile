#BUILDING
FROM golang:1.20-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLE=0 go build -o my-go-app ./cmd/web

RUN chmod +x /app/my-go-app

#RUNNING
FROM alpine:latest

RUN  mkdir /app

COPY --from=builder /app/my-go-app /app

CMD [ "/app/my-go-app" ]
