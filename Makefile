.PHONY: build test run docker-build docker-run

build:
	go build -o bin/app cmd/app/main.go

test:
	go test ./...
	
run: build
	./bin/app

docker-build:
	docker build -t mydoctor-server .

docker-run:
	docker run -p 8080:8080 mydoctor-server
