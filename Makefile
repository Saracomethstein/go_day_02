TEST_DIR=test_dir
BUILD_DIR=build

FIND_APP=$(BUILD_DIR)/find

.PHONY: test clean build generate_build_dir generate_test_dir 

build: generate_build_dir
	@go build -o $(FIND_APP) cmd/find/main.go

test: generate_dir
	@go test -v ./pkg/find

generate_build_dir:
	@mkdir $(BUILD_DIR)

generate_test_dir:
	@sh tests/init_test_dir.sh

clean:
	@rm -rf $(TEST_DIR)
	@rm -rf $(BUILD_DIR)
