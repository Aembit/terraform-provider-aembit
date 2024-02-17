default: testacc

install:
	go install .

# Run acceptance tests
.PHONY: testacc
testacc: install
	cd internal/provider
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
