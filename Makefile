.PHONY: test
test:
	docker build -t test --target=src .
	docker run --rm test go test -v -cover ./... --tags=integration

.PHONY: cover
cover:
	docker build -t test --target=src .
	docker run --rm -v `pwd`:/test-results test go test -v -covermode count -coverprofile /test-results/coverage.coverprofile ./... --tags=integration

.PHONY: coveralls
coveralls: cover
	gover
	goveralls -coverprofile gover.coverprofile -service travis-ci -repotoken $(COVERALLS_TOKEN)