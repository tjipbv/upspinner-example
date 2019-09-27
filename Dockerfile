FROM golang:1.13-alpine as builder
RUN apk --no-cache add curl git gcc libc-dev

WORKDIR /src
COPY ./ /src
RUN go test ./... -v
RUN GOOS=linux go build -a -o app main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder /src/app .
EXPOSE 3000
CMD ["/app/app"]
