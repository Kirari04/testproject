
build:
	@go build -o main ./main.go

run:
	@go run ./main.go

build-docker:
	@docker build -t testproject:latest . --load

run-docker:
	@docker run --name testproject --rm -p 8080:8080 testproject:latest