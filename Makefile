.PHONY: test
test:
	docker build -t test --target=src .
	docker run --rm test go test -v -cover ./... --tags=integration