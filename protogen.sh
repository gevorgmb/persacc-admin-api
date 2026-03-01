protoc \
  --proto_path=./proto \
  --go_out=./api/v1/admin \
  --go_opt=paths=source_relative \
  --go-grpc_out=./api/v1/admin \
  --go-grpc_opt=paths=source_relative \
  ./proto/*.proto
