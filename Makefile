start:
	go run cmd/app/main.go

fmt:
	go fmt ./...

wire-gen:
	wire ./internal/...

swag-gen:
	swag fmt ./internal/adapters/rest/...
	swag init -g cmd/app/main.go -o ./docs

generate: wire-gen swag-gen
