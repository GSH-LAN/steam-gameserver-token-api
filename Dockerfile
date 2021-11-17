####### BUILD ENVIRONMENT #######
FROM golang:1.17.2-alpine as builder

# RUN apk add --no-cache --virtual .build-deps \
# 	alpine-sdk \
# 	cmake

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY ./src ./src

RUN CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go build -o app ./src

####### PROD IMAGE #######
FROM alpine:3.14
RUN addgroup -g 1000 -S go && \
    adduser -u 1000 -S web -G go && \
    apk add --no-cache ca-certificates tzdata

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /build/app ./app

USER web

EXPOSE 8080

CMD ["./app"]
