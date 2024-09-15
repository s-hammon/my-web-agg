build:
	@go build -o out 

run:
	@./out

clean:
	@rm out
	@go mod tidy

out: build run