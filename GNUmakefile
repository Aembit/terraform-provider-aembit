default: test_coverage
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
# Must use docker pull first to make sure we actually have the latest
	docker pull golangci/golangci-lint:latest
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run -v

install:
	go install -a -ldflags "-X main.version=1.0.0" .

# Run acceptance tests
.PHONY: testacc unittest test_coverage
testacc: .check-env-vars install
	cd internal/provider
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 10m -parallel 4 -coverprofile coverage_at.out
	printf 'REMINDER: To run specific tests, use \e[36mTESTARGS="-run REGEX"\e[0m\n'

unittest: .check-env-vars install
	cd internal/provider
	go test ./... -v $(TESTARGS) -timeout 10m -parallel 4 -coverprofile coverage_ut.out
	printf 'REMINDER: To run specific tests, use \e[36mTESTARGS="-run REGEX"\e[0m\n'

test_coverage: testacc unittest
	echo 'mode: set' > coverage.out
	tail -q -n +2 coverage_at.out coverage_ut.out >> coverage.out
	go tool cover -html coverage.out -o coverage.html

# Locally create a build for local/qa testing using GoReleaser
#	Reference: https://developer.hashicorp.com/terraform/registry/providers/publishing#using-goreleaser-locally
build: test_coverage
	go install github.com/goreleaser/goreleaser@latest
	goreleaser build --snapshot --clean

# Generate updated docs bundle
#   Individual docs files can be tested using the Terraform Registry Docs Preview
#		available at https://registry.terraform.io/tools/doc-preview
docs:
	go generate ./...
