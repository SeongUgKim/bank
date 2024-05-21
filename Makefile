createdb:
	docker exec -it postgres_container createdb --username=postgres --owner=postgres simple_bank
migrateup:
	migrate -path db/migration -database "postgresql://postgres:mysecretpassword@localhost:5433/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://postgres:mysecretpassword@localhost:5433/simple_bank?sslmode=disable" -verbose down
server:
	go run main.go
.PHONY: createdb migrateup migratedown server