.PHONY: docker
docker:
	@rm webook || true
	@go build -o webook .
	@docker rmi -f curry/webook:v0.0.1
	@docker build -t curry/webook:v0.0.1 .