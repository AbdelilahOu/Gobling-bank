postgres:
	docker run --name postgres-database-dev -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres:15

createdb: 
	docker exec -it postgres-database-dev createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres-database-dev dropdb simple_bank

migrationup:
	migrate -path db/migration -database "postgres://root:mysecretpassword@localhost:5432/simple_bank?sslmode=disable" --verbose up

migrationdown:
	migrate -path db/migration -database "postgres://root:mysecretpassword@localhost:5432/simple_bank?sslmode=disable" --verbose down

generategofunc: 
	docker run --rm -v "%cd%:/src" -w /src sqlc/sqlc generate

.PHONY: postgres createdb dropdb migrationup migrationdown