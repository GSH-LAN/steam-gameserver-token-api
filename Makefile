.PHONY: \
	build \
	test
run:
	go run ./src/main.go

build:
	go build ./...

test:
	go test -v -count=1 -race ./...

docker:                                                                                                       .
	docker buildx build -t ghcr.io/gsh-lan/steam-gameserver-token-api:latest . --platform=linux/amd64

docker-arm:
	docker buildx build -t ghcr.io/gsh-lan/steam-gameserver-token-api:latest . --platform=linux/arm64

dockerx:
	docker buildx create --name steam-gameserver-token-api-builder --use --bootstrap
	docker buildx build -t ghcr.io/gsh-lan/steam-gameserver-token-api:latest --platform=linux/arm64,linux/amd64 .
	docker buildx rm steam-gameserver-token-api-builder

dockerx-builder:
	docker buildx create --name steam-gameserver-token-api-builder --use --bootstrap
