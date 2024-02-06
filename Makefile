.PHONY: docker
docker:
	@rm webook || true
	@go build -o webook .
	@docker rmi -f evantheboy/webook-live:v0.0.1
	@docker build -t evantheboy/webook-live:v0.0.1 .