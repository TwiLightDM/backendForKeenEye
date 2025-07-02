COVERAGE_FILE=coverage.out

test:
	go test -v ./...

cover:
	go test -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -func=$(COVERAGE_FILE)

cover-html:
	go tool cover -html=$(COVERAGE_FILE) -o coverage.html

swagger:
	swag init
	go run .

back:
	docker compose up -d --build back-go
