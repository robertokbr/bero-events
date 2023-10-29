SERVER_SOURCE = src/infra/cmd/main.go

# Targets
.PHONY: server test reset build

run:
	@go run $(SERVER_SOURCE)

build:
	@go build -o bin/api $(SERVER_SOURCE)
