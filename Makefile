BINARY_NAME=myapp
BUILD_DIR=./cmd/web/bin
PORT=8000
DSN="host=host.docker.internal port=1234 user=postgres password=password dbname=concurrency sslmode=disable timezone=UTC connect_timeout=5"
REDIS="host.docker.internal:6379"

## build: Build binary
build:
	@echo "Building..."
	@env go build -v -o $(BUILD_DIR)/$(BINARY_NAME).exe ./cmd/web
	@echo "Built!"

## run: builds and runs the application
run: build
	@echo "Starting..."
	@set DSN=$(DSN) && set REDIS=$(REDIS) && $(BUILD_DIR)/$(BINARY_NAME).exe
	@echo "Started!"

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@go clean
	@rm ./cmd/web/bin/${BINARY_NAME}
	@echo "Cleaned!"

## start: an alias to run
start: run

## stop: stops the running application
stop:
	@echo "Stopping..."
	@for /f "tokens=5" %%a in ('netstat -ano ^| findstr :$(PORT) ^| findstr LISTENING') do taskkill /F /PID %%a
	@echo "Stopped!"

docker-up:
	@echo "Starting Docker services..."
	@docker-compose up -d

## docker-down: stop docker services
docker-down:
	@echo "Stopping Docker services..."
	@docker-compose down
## restart: stops and starts the application
restart: docker-down stop docker-up start

## test: runs all tests
test:
	go test -v ./...