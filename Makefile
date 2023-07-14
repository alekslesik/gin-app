PROJECT := gin-app
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
######### SEMGREP #########
# run semgrep tests
.PHONY: semgrep
semgrep: $(ALLGO) $(RULES)
	. venv/bin/activate && semgrep --config rules/ --metrics=off

# run semgrep tests
# RULES := $(wildcard rules/*.yaml)
# .PHONY: semgrep
# semgrep: $(ALLGO) $(RULES)
# 	semgrep --config rules/ --metrics=off
