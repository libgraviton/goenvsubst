FROM golang:1.13-alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-s -w -extldflags "-static"' -o goenvsubst .

FROM alpine:latest
COPY --from=builder /build/goenvsubst /app/
WORKDIR /app
CMD ["./goenvsubst"]
