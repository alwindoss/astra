FROM golang:1.17 AS builder
WORKDIR /go/src/app
COPY . .
RUN make setup
RUN make docker

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/app/bin/astra .
CMD [ "./astra" ]