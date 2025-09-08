all: clean test itest run

build:
	@echo 'Building...'
	@go build -o limitr .

clean:
	@echo 'Cleaning...'
	@rm limitr

run: build
	@echo 'Starting...'
	@./limitr

test:
	@echo "Testing..."
	@go test ./... -v

itest:
	@echo "Testing [integration]..."
	@go test ./internal/database -v
