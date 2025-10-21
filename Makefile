all: clean test run

build:
	@echo 'Building...'
	@go build -o limitr .

clean:
	@echo 'Cleaning...'
	@rm limitr

test:
	@echo 'Testing...'
	@go test ./... -v

run: build
	@echo 'Starting...'
	@./limitr

docker:
	@echo 'Starting Docker...'
	@docker compose up --build
	@echo 'Stopping Docker...'
	@docker compose down

docker-rebuild:
	@echo 'Forcing a fresh Docker build...'
	@docker compose build --no-cache app
	@docker compose up --force-recreate
	@echo 'Stopping Docker...'
	@docker compose down
