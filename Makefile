build:
	go build cmd/main.go

run:
	go run cmd/main.go

unit-test:
	go test -v ./pkg/...

integration-test:
	go test -v ./test/...

test: unit-test integration-test

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		pkg/infrastructure/grpc/proto/ingredients_balancer.proto
