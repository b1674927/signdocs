# Variables
APP_NAME = signdocs
SOURCE = ./cmd/main.go
BUILD_DIR = ./build

# Default target
all: build

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(SOURCE)
	@echo "Build complete. Executable is at $(BUILD_DIR)/$(APP_NAME)"

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete."

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	@$(BUILD_DIR)/$(APP_NAME)

.PHONY: all build clean run
