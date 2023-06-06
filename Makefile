postgres:
	docker run --name --network bank-network postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:15-alpine

createdb:
	docker exec -it postgres15 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres15 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:fzRSwkNUzlTUZsXVEajH@simple-bank.cofr1fxqo7ie.ca-central-1.rds.amazonaws.com:5432/simple_bank" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:fzRSwkNUzlTUZsXVEajH@simple-bank.cofr1fxqo7ie.ca-central-1.rds.amazonaws.com:5432/simple_bank" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://root:fzRSwkNUzlTUZsXVEajH@simple-bank.cofr1fxqo7ie.ca-central-1.rds.amazonaws.com:5432/simple_bank" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://root:fzRSwkNUzlTUZsXVEajH@simple-bank.cofr1fxqo7ie.ca-central-1.rds.amazonaws.com:5432/simple_bank" -verbose down 1

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -build_flags=--mod=mod -package mockdb -destination db/mock/store.go github.com/techschool/simplebank/db/sqlc Store

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto

evans:
	evans --host localhost --port 8081 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7-alpine


.PHONY: createdb dropdb migrateup migratedown migrateup1 migratedown1 new_migration test server mock proto evans redis
