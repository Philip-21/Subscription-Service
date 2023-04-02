BINARY_NAME=my-go-app.exe
DSN="host=localhost port=5432 user=postgres password=philippians dbname=plans sslmode=disable timezone=UTC connect_timeout=5"
REDIS="127.0.0.1:6379"

up_build :
	@echo "Stopping Docker Images if Running...."
	docker-compose down 
	@echo "Building Docker Images were necessary"
	docker-compose up --build 
up :  
	@echo "Starting Docker Images"
	docker-compose up -d
	@echo "Docker Images Started"
      
## build: Build binary
build:
	@echo "Building..."
	@go build -o ${BINARY_NAME} ./cmd/web
	@echo "Built!"

## run: builds and runs the application
run: build
	@echo Starting...
	powershell $$env:DSN='${DSN}'; $$env:REDIS='${REDIS}'; ./${BINARY_NAME}
	@echo Started!

## clean: runs go clean and deletes binaries
clean:
	@echo "Cleaning..."
	@go clean
	@rm ${BINARY_NAME}
	@echo "Cleaned!"

## start: an alias to run
start: run

## stop: stops the running application
stop:
	@echo "Stopping..."
	@-pkill -SIGTERM -f "./${BINARY_NAME}"
	@echo "Stopped!"

## restart: stops and starts the application
restart: stop start

## test: runs all tests
test:
	go test -v ./...