TEST?=$$(go list ./... | grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
HOSTNAME=lhalbert.xyz
NAMESPACE=lucashalbert
NAME=utilities
BINARY=terraform-provider-${NAME}
VERSION=0.0.4
OS=linux
ARCH=amd64
OS_ARCH=${OS}_${ARCH}

default: build

build:
	@echo "==> Running Code Build..."
	@GOOS=${OS} GOARCH=${ARCH} ENBABLE_CGO=0 go build -o bin/${BINARY}_${OS_ARCH}

clean:
	@echo "==> Running cleanup methods..."
	@go clean
	@rm bin/terraform-provider-utilities_linux_amd64
	@unlink ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/${BINARY}

install: build
	@echo "==> Running Code Install..."
	@mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	@ln -s $$(pwd)/bin/${BINARY}_${OS_ARCH} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/${BINARY}

fmt:
	@echo "==> Running go fmt..."
	@gofmt -w $(GOFMT_FILES)

testall: test testcov testrace testacc

test:
	@echo "==> Running Tests..."
	go test $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4                    

testcov:
	@echo "==> Running Code Coverage Tests..."
	go test -cover $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4  

testrace:
	@echo "==> Running Code Race Condition Tests..."
	go test -race $(TEST) || exit 1                                                   
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4  

testacc:
	@echo "==> Running Acceptance Testing..."
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m   

semgrep:
	@echo "==> Running Semgrep static analysis..."
	@docker run --rm --volume "${PWD}:/src/" returntocorp/semgrep semgrep --config "p/golang" --metrics=off --verbose

