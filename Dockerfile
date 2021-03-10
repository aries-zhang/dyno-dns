FROM golang:1.15.0-alpine

WORKDIR /go/src/github.com/aries-zhang/dyno-dns/
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=0 /go/src/github.com/aries-zhang/dyno-dns/app .
CMD ["./app"]