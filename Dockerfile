FROM golang:alpine as builder

WORKDIR /app 

COPY . /app

RUN go build -o /tmp/app

FROM alpine

COPY --from=builder /tmp/app /app

EXPOSE 8080

CMD /app
