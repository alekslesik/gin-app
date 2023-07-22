PROJECT := gin-app
CMDS := importer
ALLGO := $(wildcard *.go */*.go cmd/*/*.go)
ALLHTML := $(wildcard templates/*/*.html)

.DELETE_ON_ERROR:

.PHONY: all
all: lint check $(PROJECT) $(CMDS)

.PHONY: lint
lint: .lint

.lint: $(ALLGO)
	golangci-lint run --timeout=180s --skip-dirs=rules
	@touch $@

# create .coverage dir
.coverage:
	@mkdir -p .coverage

.PHONY: check
check: .coverage ./.coverage/$(PROJECT).out

./.coverage/$(PROJECT).out: $(ALLGO) $(ALLHTML) Makefile
	go test $(TESTFLAGS) -coverprofile=./.coverage/$(PROJECT).out .

.PHONY: cover
# When running manually, capture just the total percentage (and
# beautify it slightly because the tool output is usually too wide).
cover: .coverage ./.coverage/$(PROJECT).html
	@echo "Checking overall code coverage..."
	@go tool cover -func .coverage/$(PROJECT).out | sed -n -e '/^total/s/:.*statements)[^0-9]*/: /p'

./.coverage/$(PROJECT).html: ./.coverage/$(PROJECT).out
	go tool cover -html=./.coverage/$(PROJECT).out -o ./.coverage/$(PROJECT).html

# XXX *.go isn't quite right here -- it will rebuild when tests are
# touched, but it's good enough.
$(PROJECT): $(ALLGO)
	go build .

$(CMDS): $(ALLGO)
	for cmd in $(CMDS); do go build ./cmd/$$cmd; done

# XXX This only works if go-junit-report is installed. It's not part of go.mod
# because I don't want to force a dependency, but it is part of the ci docker
# image.
report.xml: $(ALLGO) Makefile
	go test $(TESTFLAGS) -v . 2>&1 | go-junit-report > $@
	go tool cover -func .coverage/$(PROJECT).out

RULES := $(wildcard rules/*.yaml)

#=====================================#
# SEMGREP #
#=====================================#

# run semgrep tests
.PHONY: semgrep
semgrep: $(ALLGO) $(RULES)
	. venv/bin/activate && semgrep --config rules/ --metrics=off

#=====================================#
# HELPERS #
#=====================================#

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

#=====================================#
# DEVELOPMENT #
#=====================================#

## run: run the cmd/app application
.PHONY: run
run: audit
	go run .

#=====================================#
# QUALITY CONTROL #
#=====================================#

## audit: tidy dependencies and format, vet and test all code

## go fmt ./... : command to format all .go files in the project directory, according to the Go standard.
## go vet ./... : runs a variety of analyzers which carry out static analysis of your code and warn you
## go test -race -vet=off ./... : command to run all tests in the project directory
## staticcheck tool : to carry out some additional static analysis checks.
.PHONY: audit
audit: vendor
	@echo 'Formatting code...'
	go fmt ./...
	@echo 'Vetting code...'
	go vet ./...
	# staticcheck ./...
	# @echo 'Running tests...'
	# go test -race -vet=off ./...

## go mod tidy : prune any unused dependencies from the go.mod and go.sum files, and add any missing dependencies
## go mod verify : check that the dependencies on your computer (located in your module cache located at $GOPATH/pkg/mod)
## haven’t been changed since they were downloaded and that they match the cryptographic hashes in your go.sum file
## go mod vendor: copy the necessary source code from your module cache into a new vendor directory in your project root
.PHONY: vendor
vendor:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	# @echo 'Vendoring dependencies...'
	# go mod vendor

## import csv table "./importer -db=gin-app.db -csv=goodreads_library_export.csv"
.PHONY: csv.import
csv.import:
	./importer -db=gin-app.db -csv=goodreads_library_export.csv
