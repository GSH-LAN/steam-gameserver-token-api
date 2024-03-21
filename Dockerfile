FROM golang:1.21.7-alpine as builder

RUN apk --no-cache add ca-certificates

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY ./src ./src

RUN CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o app ./src

FROM scratch

# copy the ca-certificate.crt from the build stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /build/app ./app

EXPOSE 8080

CMD ["./app"]
