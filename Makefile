CURRENT_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

containerup:
	docker run --name postgres-database-dev -e POSTGRES_USER=root -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres:15

containerdown:
	docker stop postgres-database-dev
	docker rm --force postgres-database-dev

createdb: 
	docker exec -it postgres-database-dev createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres-database-dev dropdb simple_bank

migrationup:
	migrate -path db/migration -database "postgres://root:mysecretpassword@localhost:5432/simple_bank?sslmode=disable" --verbose up

migrationdown:
	migrate -path db/migration -database "postgres://root:mysecretpassword@localhost:5432/simple_bank?sslmode=disable" --verbose down

sqlc: 
	docker run --rm -v ${CURRENT_DIR}:/src -w /src sqlc/sqlc generate

test:
	go test -v -cover ./...

respawn:
	make containerdown
	make containerup
	timeout 3
	make createdb
	make migrationup
	make sqlc
	make test

.PHONY: postgres createdb dropdb migrationup migrationdown sqlc test respawn