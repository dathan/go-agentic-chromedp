#!/bin/bash

# install.sh – Set up dependencies for the agentic Chrome Go program on macOS.
# This script installs Homebrew if it isn't present, then uses Homebrew to
# install Go, Google Chrome, and Ollama. It is designed for Apple Silicon
# Macs (M1/M2/M3), as described in the project README. If you already
# have these tools installed, the script skips those steps.

set -e

# Confirm we are running on macOS. The script exits if it is not.
if [[ "$OSTYPE" != "darwin"* ]]; then
  echo "This installer is intended for macOS. Exiting." >&2
  exit 1
fi

# Function to check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

echo "Checking for Homebrew..."
if ! command_exists brew; then
  echo "Homebrew not found. Installing Homebrew..."
  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
else
  echo "Homebrew is already installed."
fi

echo "Updating Homebrew..."
brew update

# Install Go if missing
if ! command_exists go; then
  echo "Installing Go..."
  brew install go
else
  echo "Go is already installed."
fi

# Install Google Chrome if missing
if ! command_exists google-chrome && ! command_exists "Google Chrome"; then
  echo "Installing Google Chrome..."
  brew install --cask google-chrome
else
  echo "Google Chrome is already installed."
fi

# Install Ollama if missing
if ! command_exists ollama; then
  echo "Installing Ollama..."
  brew install ollama
else
  echo "Ollama is already installed."
fi

echo "Environment setup complete."
