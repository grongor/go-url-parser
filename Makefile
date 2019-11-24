BIN = ${PWD}/bin

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/url-parser

.PHONY: check
check: cs vet staticcheck test

.PHONY: cs
cs: tools
	diff=$$($(BIN)/goimports -d . ); test -z "$$diff" || (echo "$$diff" && exit 1)

.PHONY: cs-fix
cs-fix: format

.PHONY: format
format: tools
	$(BIN)/goimports -w .

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	timeout 10 go test
	timeout 60 go test --race
	timeout 120 go test --count 1000
	timeout 10 cmd/url-parser/main_test.sh

.PHONY: staticcheck
staticcheck: tools
	$(BIN)/staticcheck ./...

.PHONY: tools
tools: bin

bin: export GOBIN = $(BIN)
bin:
	go install github.com/golang/mock/mockgen
	go install golang.org/x/lint/golint
	go install golang.org/x/tools/cmd/goimports
	go install honnef.co/go/tools/cmd/staticcheck
