DB_URL=postgresql://root:secret@localhost:5433/simple_bank?sslmode=disable

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/SeongUgKim/simplebank/db/sqlc Store
.PHONY: createdb, dropdb migrateup migratedown sqlc test server mock
