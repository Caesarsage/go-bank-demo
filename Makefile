build:
	@go build -o bin/gobank
run: build
	@./bin/gobank
test:
	@go-test -v ./..
seed:
	@./bin/go-bank --seed true
