FROM golang:1.21.1 as builder

WORKDIR /go/src/builder
COPY . .
RUN make build

FROM alpine:3
COPY --from=builder /usr/bin/tee /bin/tee
COPY --from=builder /go/src/builder/build/go-boilerplate /app/go-boilerplate
WORKDIR /app
EXPOSE 3001
CMD ["./go-boilerplate"]
