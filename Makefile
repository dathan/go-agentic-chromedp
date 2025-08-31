## Makefile for the agentic Chrome Go program
#
# This Makefile defines common tasks for building, running and installing
# dependencies on a MacBook Pro with an M1/M2/M3 chip (Apple Silicon). It
# assumes that your environment includes the install script located in
# scripts/install.sh. Use `make install` to prepare your Mac for the
# program, `make build` to build the binary, and `make run` to build
# and execute the program. Additional targets can be added as needed.

# Name of the resulting binary when building
BIN_NAME := agentic-chromedp

.PHONY: build run install clean

# Build the Go program into a binary in the current directory
build:
	@echo "Building $(BIN_NAME)..."
	go build -o $(BIN_NAME) main.go

# Run the program; depends on build to ensure the binary exists
run: build
	@echo "Running $(BIN_NAME)..."
	./$(BIN_NAME)

# Install dependencies using the provided install script
install:
	@echo "Installing dependencies via scripts/install.sh..."
	chmod +x scripts/install.sh
	./scripts/install.sh

# Remove the built binary
clean:
	@echo "Cleaning up build artifacts..."
	rm -f $(BIN_NAME)
