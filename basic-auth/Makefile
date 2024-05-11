
start:
	go run ./cmd/ --debug

build-docker:
	docker \
		buildx \
		build \
		-t "eldius/basic-auth:dev" \
			.

tidy:
	go mod tidy

test:
	go test ./... -cover
