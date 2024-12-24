TEST_DIR=test_dir
BUILD_DIR=build

FIND_APP=$(BUILD_DIR)/find
WC_APP=$(BUILD_DIR)/wc

.PHONY: test clean build generate_build_dir generate_test_dir 

all: build

build: generate_build_dir
	@go build -o $(FIND_APP) cmd/find/main.go
	@go build -o $(WC_APP) cmd/wc/main.go

test: generate_dir
	@go test -v ./pkg/find

generate_build_dir:
	@mkdir $(BUILD_DIR)

generate_test_dir:
	@sh tests/init_test_dir.sh

clean:
	@rm -rf $(TEST_DIR)
	@rm -rf $(BUILD_DIR)
