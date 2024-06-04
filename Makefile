
build:
	@go build -o main ./main.go

run:
	@go run ./main.go

build-docker:
	@docker build -t kirari04/haproxymngr:latest . --load

run-docker:
	@docker run --name haproxymngr --rm -p 8080:8080 kirari04/haproxymngr:latest

release-docker:
	@docker build -t kirari04/haproxymngr:latest . --push