FROM golang:alpine AS builder
RUN apk add build-base
WORKDIR /go/src/app
COPY ./ /go/src/app/
RUN go build -o main cmd/main.go
FROM alpine
WORKDIR /app
COPY --from=builder /go/src/app/ /app/
CMD [ "./main" ]