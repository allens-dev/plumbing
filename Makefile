.PHONY: test lint

test:
	go test -timeout 30s -coverprofile=test-coverage.out ./...

lint:
	golangci-lint run -c ./.golangci.yml ./...