default: testacc

install:
	go mod tidy
	go install github.com/goreleaser/goreleaser@latest
	go get aembit.io/aembit
	go install -a -ldflags "-X main.version=1.0.0" .

# Run acceptance tests
.PHONY: testacc
testacc: install
	cd internal/provider
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 10m

# Locally create a build for local/qa testing using GoReleaser
#	Reference: https://developer.hashicorp.com/terraform/registry/providers/publishing#using-goreleaser-locally
build: testacc
	goreleaser build --snapshot --clean

release: testacc
	goreleaser release