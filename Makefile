auth_protoc:
	protoc -I ./protos/proto \
      --go_out ./protos/gen --go_opt paths=source_relative \
      --go-grpc_out ./protos/gen --go-grpc_opt paths=source_relative \
      --grpc-gateway_out ./protos/gen --grpc-gateway_opt paths=source_relative \
      --openapiv2_out ./protos/gen --openapiv2_opt use_go_templates=true \
      ./protos/proto/auth.proto

run:
	docker compose -f docker-compose.yml up

down:
	docker compose -f docker-compose.yml down

build:
	docker build -t grpc-auth-service .