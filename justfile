python := "python3"
pip := "pip"
go := "go"

[private]
default:
  @just --list

# Generate protobuf files for the Python server
[private]
create-python-protobuf:
	{{ python }} -m grpc_tools.protoc -I. --python_out=./internal/python/ --grpc_python_out=./internal/python/ ./shareProfileAllocator.proto

# Generate protobuf files for the Go server
[private]
create-go-protobuf:
	@set -e
	protoc --go_out=. --go-grpc_out=. ./shareProfileAllocator.proto

# Create all the protobuf files for the server
create-protobuf-files: create-go-protobuf create-python-protobuf

# Delete all generated protobuf files
delete-protobuf-files:
	rm -r ./internal/grpc/generated/go/
	rm ./internal/python/shareProfileAllocator_pb2.py
	rm ./internal/python/shareProfileAllocator_pb2_grpc.py

# Install all project dependencies
install-deps:
	{{ pip }} install -r requirements.txt
	{{ go }} install
	@echo 'I recommend installing jq'

# Run the python finance share server independently
run-share-server:
	{{ python }} ./internal/python/finance_data_manager.py

# Run the html server and the finance server
run finance-server="true":
    go run ./main.go --finance-server={{finance-server}} | jq
