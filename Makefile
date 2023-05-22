generate:
	go generate ./...

gen: generate

dev:
	go run cmd/.main.go

