protoc:
	protoc -I protos/proto protos/proto/auth.proto --go_out=./protos/gen --go_opt=paths=source_relative --go-grpc_out=./protos/gen/ --go-grpc_opt=paths=source_relative

run:
	docker compose --env-file .env -f docker-compose.yml up

down:
	docker compose -f docker-compose.yml down