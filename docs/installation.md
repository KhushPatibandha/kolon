# Installation Guide for [Your Language]

Welcome to the installation guide for **Kolon**! Follow the steps below to install and set up the language on your system.

## Prerequisites

Before installing **Kolon**, ensure you have the following:

- **Go Programming Language** (version 1.20 or higher) installed on your system.
  - [Download and install Go](https://golang.org/dl/)

## Installation Methods

You can install **Kolon** using one of the two methods below:

### Method 1: Install Prebuilt Binary

1. Visit the [Releases Page](https://github.com/your-repo/your-language/releases).
2. Download the prebuilt binary for your operating system:
   - **Linux**: `kolon_x86_64`
3. Make the binary executable (Linux/MacOS):
   ```bash
   chmod +x kolon_x86_64
   ```
4. Move the binary to a directory in your `PATH` (e.g., `/usr/local/bin`):
   ```bash
   sudo mv kolon_x86_64 /usr/local/bin/kolon
   ```
5. Verify the installation:
   ```bash
   kolon --version
   ```

### Method 2: Build from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/KhushPatibandha/Kolon.git
   cd Kolon
   ```
2. Build the binary using Go:
   ```bash
    go build -o kolon cmd/main.go
   ```
3. Move the binary to a directory in your `PATH` (e.g., `/usr/local/bin`):
   ```bash
   sudo mv kolon /usr/local/bin/
   ```
4. Verify the installation:
   ```bash
   kolon --version
   ```

For more information, refer to the [documentation](https://github.com/your-repo/your-language/wiki).
