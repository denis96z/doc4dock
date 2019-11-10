CMD_DIR := $(PWD)/cmd
PKG_DIR := $(PWD)/pkg

MAIN_PATH := $(CMD_DIR)/doc4dock/main.go

BUILD_DIR := $(PWD)/build
TESTS_DIR := $(PWD)/tests
TOOLS_DIR := $(PWD)/tools

BIN_FILE := $(BUILD_DIR)/doc4dock
COVERAGE_PROFILE_FILE := $(TESTS_DIR)/c.out
COVERAGE_REPORT_FILE := $(TESTS_DIR)/c.html

BUILD_FLAGS := -mod=vendor
INSTALL_FLAGS := -mod=vendor
TEST_FLAGS := -mod=vendor -v -race -coverprofile='$(COVERAGE_PROFILE_FILE)' -covermode=atomic
COVER_FLAGS := -html='$(COVERAGE_PROFILE_FILE)' -o '$(COVERAGE_REPORT_FILE)'

.PHONY: fmt
fmt:
	go fmt $(PWD)/...

.PHONY: dep
dep:
	go mod tidy && go mod vendor && go mod verify

.PHONY: lint
lint: bootstrap fmt
	'$(TOOLS_DIR)/golangci-lint' run

.PHONY: build
build: dep
	GOBIN='$(BUILD_DIR)' go build $(BUILD_FLAGS) -o $(BIN_FILE) $(MAIN_PATH)

.PHONY: test
test: dep
	go test $(TEST_FLAGS) $(PWD)/...
	GOFLAGS='-mod=vendor' go tool cover $(COVER_FLAGS)

.PHONY: clean
clean:
	rm -f '$(BIN_FILE)'
	rm -f '$(COVERAGE_PROFILE_FILE)'
	rm -f '$(COVERAGE_REPORT_FILE)'

.PHONY: mkdirs
mkdirs:
	mkdir -p '$(BUILD_DIR)'
	mkdir -p '$(TESTS_DIR)'
	mkdir -p '$(TOOLS_DIR)'

.PHONY: bootstrap
bootstrap: \
	install_golangci-lint

.PHONY: install_golangci-lint
install_golangci-lint: mkdirs dep
	GOBIN='$(TOOLS_DIR)' go install $(INSTALL_FLAGS) \
		$(PWD)/vendor/github.com/golangci/golangci-lint/cmd/golangci-lint
