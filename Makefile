.PHONY: test
test:
	docker build -t test --target=src .
	docker run --rm test go test -mod vendor -v -cover ./... --tags=integration