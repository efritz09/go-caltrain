
.PHONY: test
test:
	go test ./... -race -cover -coverprofile=c.out -count=1

.PHONY: lint
lint:
	golangci-lint run
