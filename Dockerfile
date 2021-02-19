FROM golang:1.16 AS builder
WORKDIR /src
COPY . /src
RUN CGO_ENABLED=0 GOOS=linux go build -o starlink-exporter cmd/starlink-exporter/main.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /src/starlink-exporter .
CMD ["./starlink-exporter"]  
