default: testacc
.check-env-vars:
	@ if [ -z $$AEMBIT_TENANT_ID ]; then \
        echo "Environment variable AEMBIT_TENANT_ID not set"; \
        exit 1; \
	fi
	@ if [ -z $$AEMBIT_STACK_DOMAIN ]; then \
        echo "Environment variable AEMBIT_STACK_DOMAIN not set"; \
        exit 1; \
	fi
	@ if [ -z $$AEMBIT_TOKEN ]; then \
        echo "Environment variable AEMBIT_TOKEN not set"; \
        exit 1; \
	fi

# Run the GitHub CI Linters locally
lint: 
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run -v

install:
	go install -a -ldflags "-X main.version=1.0.0" .

# Run acceptance tests
.PHONY: testacc
testacc: .check-env-vars install
	cd internal/provider
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 10m -coverprofile coverage.out
	go tool cover -html coverage.out -o coverage.html

# Locally create a build for local/qa testing using GoReleaser
#	Reference: https://developer.hashicorp.com/terraform/registry/providers/publishing#using-goreleaser-locally
build: testacc
	go install github.com/goreleaser/goreleaser@latest
	goreleaser build --snapshot --clean

docs:
	go generate ./...