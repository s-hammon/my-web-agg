GOOSE_DIR="sql/schema"
CONN_STR="postgres://postgres:postgres@localhost:5432/blogator"
build:
	@go build -o out 

run:
	@./out

clean:
	@rm out
	@go mod tidy

up:
	@goose -dir ${GOOSE_DIR} postgres ${CONN_STR} up

down:
	@goose -dir ${GOOSE_DIR} postgres ${CONN_STR} down

models:
	sqlc generate

out: build run

upgrade: up models