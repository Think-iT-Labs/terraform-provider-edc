default: testacc

# Run acceptance tests
.PHONY: testacc unit-test
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

unit-test:
	go test -coverprofile=coverage.out ./...

pre-commit:
	pre-commit install --install-hooks -t pre-commit -t commit-msg
docs-gen:
	go generate ./...
