.PHONY: test
test:
	docker build -t test --target=src .
	docker run --rm test go test -v -cover ./... --tags=integration

.PHONY: cover
cover:
	docker build -t test --target=src .
	docker run --rm -v `pwd`:/test-results test go test -v -covermode count -coverpkg ./... -coverprofile /test-results/coverage.coverprofile ./... --tags=integration

.PHONY: coveralls
coveralls: cover
	goveralls -coverprofile coverage.coverprofile -service travis-ci -repotoken $(COVERALLS_TOKEN)