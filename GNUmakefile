default: testacc

install:
	go get aembit.io/aembit
	go install .

# Run acceptance tests
.PHONY: testacc
testacc: install
	cd internal/provider
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m
