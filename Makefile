.PHONY: gen test build harness

gen:
	go run ./cmd/qapi-gen --schema-dir ./schemas/qemu-5.0.0/qapi --out-dir ./internal/gen/v5 --version 5.0.0 --package gen_v5
	go run ./cmd/qapi-gen --schema-dir ./schemas/qemu-10.2.2/qapi --out-dir ./internal/gen/v10 --version 10.2.2 --package gen_v10

test:
	go test ./...

build:
	go build -o bin/qapi-gen ./cmd/qapi-gen
	go build -o bin/qvm ./cmd/qvm

harness:
	go run ./harness/cmd/run
