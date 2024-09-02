dockertag = kirari04/haproxymngr:latest
build:
	@go build -o main ./main.go

run:
	@go build -o main ./main.go && ./main serve

dev:
	@air serve

build-docker:
	@docker build -t $(dockertag) . --load

run-docker:
	@docker run --name haproxymngr --rm --network host -v haproxymngr:/app/data -e ADDR=0.0.0.0:8080 $(dockertag)

release-docker:
	@docker build -t $(dockertag) . --push