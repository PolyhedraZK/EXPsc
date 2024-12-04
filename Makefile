
BINARY_NAME=sc
CMD_DIR=./cmd
CORE_DIR=./core
OUTPUT_DIR=../


build:
	cd $(CMD_DIR) && go build -o $(OUTPUT_DIR)$(BINARY_NAME)

tidy:
	cd $(CMD_DIR) && go mod tidy
	cd $(CORE_DIR) && go mod tidy