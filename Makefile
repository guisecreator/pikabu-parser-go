generate:
	go generate ./...

gen: generate

dev:
	go run .main.go

