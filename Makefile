dockertag = kirari04/haproxymngr:latest
build:
	@go build -o main ./main.go

run:
	@go run ./main.go

build-docker:
	@docker build -t $(dockertag) . --load

run-docker:
	@docker run --name haproxymngr --rm -p 8080:8080 $(dockertag)

release-docker:
	@docker build -t $(dockertag) . --push