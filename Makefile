all: clean test run

build:
	@echo 'Building...'
	@go build -o limitr .

clean:
	@echo 'Cleaning...'
	@rm limitr

run: build
	@echo 'Starting...'
	@./limitr

docker-run:
	@echo 'Starting Docker...'
	@docker compose up --build

docker-rebuild:
	@echo 'Forcing a fresh Docker build...'
	@docker compose build --no-cache app
	@docker compose up --force-recreate

test:
	@echo 'Testing...'
	@go test ./... -v
