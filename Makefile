run:
	go run server.go

lint:
	golangci-lint run

generate:
	go generate ./graph
