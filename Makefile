TEST_DIR=test_dir
BUILD_DIR=build

FIND_APP=$(BUILD_DIR)/find
WC_APP=$(BUILD_DIR)/wc
XARGS_APP=$(BUILD_DIR)/xargs
ROTATE_APP=$(BUILD_DIR)/rotate

.PHONY: test clean build generate_build_dir generate_test_dir 

all: build

build: generate_build_dir
	@go build -o $(FIND_APP) cmd/find/main.go
	@go build -o $(WC_APP) cmd/wc/main.go
	@go build -o $(XARGS_APP) cmd/xargs/main.go
	@go build -o $(ROTATE_APP) cmd/rotate/main.go

test: generate_test_dir build
	@echo "Run tests for FIND"
	@go test ./pkg/find
	@echo "Run tests for WC"
	@go test ./pkg/wc
	@echo "Run test for Xargs"
	@go test ./pkg/xargs
	@echo "Run test for Rotate"
	@go test ./pkg/rotate
	@rm -rf $(TEST_DIR) $(BUILD_DIR)

generate_build_dir:
	@mkdir -p $(BUILD_DIR)

generate_test_dir:
	@sh tests/scripts/init_test_dir.sh

clean:
	@rm -rf $(TEST_DIR)
	@rm -rf $(BUILD_DIR)
