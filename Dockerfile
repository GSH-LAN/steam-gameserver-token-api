FROM golang:1.24.4-alpine AS builder

RUN apk --no-cache add ca-certificates
ADD . /build
WORKDIR /build
ARG TARGETOS
ARG TARGETARCH
RUN go mod download
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -ldflags="-s -w" -o app ./src/main.go
FROM scratch

# copy the ca-certificate.crt from the build stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /build/app ./app

EXPOSE 8080

CMD ["./app"]
