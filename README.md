# Go Library Template Repository

[![Go Report Card](https://goreportcard.com/badge/github.com/csgriffis/go-library-template)](https://goreportcard.com/report/github.com/csgriffis/go-library-template)
[![Build Status](https://github.com/csgriffis/go-library-template/actions/workflows/go.yml/badge.svg)](https://github.com/csgriffis/go-library-template/actions/workflows/go.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](./LICENSE)
[![Contributors](https://img.shields.io/github/contributors/csgriffis/go-library-template)](https://github.com/csgriffis/go-library-template/graphs/contributors)
[![GoDoc](https://pkg.go.dev/badge/github.com/csgriffis/go-library-template.svg)](https://pkg.go.dev/github.com/csgriffis/go-library-template)
[![Code Coverage](https://img.shields.io/codecov/c/github/csgriffis/go-library-template)](https://codecov.io/gh/csgriffis/go-library-template)

---

Welcome to the **Go Library Template**!

This repository serves as a starting point for creating Go libraries or modules. It includes a structured project layout, essential configurations, and tools to streamline development and collaboration.

---

## Features

- **Go Modules**: Configured for dependency management.
- **Standardized Layout**: Adheres to Go's best practices for library organization.
- **CI/CD Ready**: Includes GitHub Actions workflows for testing and linting.
- **Pre-configured Tools**:
  - Linting with `golangci-lint`
  - Testing framework integration
  - Pre-commit hooks
- **Documentation Starter**: Ready-to-use outlines for common repository documents.

---

## Getting Started

### 1. Use the Template
1. Click the **"Use this template"** button on GitHub.
2. Create a new repository based on this template.

### 2. Clone the Repository
```bash
git clone https://github.com/your-username/your-new-repo.git
cd your-new-repo
```

### 3. Initialize the Module
Update the Go module name to match your new repository:
```bash
go mod edit -module github.com/your-username/your-new-repo
go mod tidy
```

### 4. Update the Badge URLs
When creating a new repository based on this template, update the the badge URLs to reflect your GitHub username and the name of your new repository.
```bash
sed -i '' 's|csgriffis/go-library-template|your-username/repo|g' README.md
```
for macOS, or if you are on Linux:
```bash
sed -i 's|csgriffis/go-library-template|your-username/repo|g' README.md
```

### 5. Update the Documents
Update the following documents for your needs:
- [ ] ARCHITECTURE.md
- [ ] CHANGELOG.md
- [ ] CODE_OF_CONDUCT.md
- [ ] CONTRIBUTING.md
- [ ] LICENSE.md
- [ ] MAINTAINERS.md
- [ ] README.md
- [ ] SECURITY.md
- [ ] SUPPORT.md

---

## Pre-commit Hooks

This repository uses [pre-commit](https://pre-commit.com/) to ensure code quality and enforce best practices before committing changes.

### 1. Install Pre-commit
First, install `pre-commit` on your local machine. You can use pip:
```bash
pip install pre-commit
```
or homebrew:
```bash
brew install pre-commit
```

### 2. Install Pre-commit Hooks
After cloning the repository, run the following command to install the pre-commit hooks:
```bash
pre-commit install
```
This will configure Git to run the hooks automatically before each commit.

### 3. Running Hooks Manually
You can run the pre-commit hooks manually on all files using:
```bash
pre-commit run --all-files
```

### 4. Included Hooks
The pre-commit configuration (.pre-commit-config.yaml) includes the following checks:
- Conventional Commits: Ensures commit messages adhere to [Conventional Commit](https://www.conventionalcommits.org/en/v1.0.0/#summary) standards.
- Code Formatting: Ensures your code is formatted properly (e.g., gofmt).
- Linting: Runs Go linters (e.g., golangci-lint) to catch potential issues.
- Static Analysis: Runs tools for detecting common issues in Go code.

### 5. Customizing Hooks
You can customize the hooks by editing the .pre-commit-config.yaml file in the repository. For details on supported hooks, visit the pre-commit hook [repository](https://github.com/pre-commit/pre-commit-hooks).

---

## Project Structure
```
golang-library-template/
├── .github/
│   └── workflows/      # CI/CD workflows
│       └── go.yml      # Go CI pipeline
├── internal/           # Non-exportable code (for internal use only)
│   └── example/
│       └── example.go  # Example internal package
├── pkg/                # Exportable code (public-facing API)
│   └── example/
│       ├── example.go       # Example public package
│       └── example_test.go  # Example test
├── go.mod                   # Module definition
├── go.sum                   # Dependencies checksum
├── .gitignore               # Files to ignore in Git
├── .golanci.yml             # golangci-lint Configuration
├── .pre-commit-config.yaml  # pre=commit hooks Configuration
├── LICENSE.md               # License
├── ARCHITECTURE.md          # Documentation
├── CHANGELOG.md
├── CODE_OF_CONDUCT.md
├── CONTRIBUTING.md
├── MAINTAINERS.md
├── SECURITY.md
├── SUPPORT.md
└── README.md
```

### Usage

#### Public API
Place all publicly accessible packages under pkg/. These are the packages users of your library will import.

#### Internal Code
Use the internal/ directory for packages that should remain private and not accessible by other modules.

### CI/CD Pipeline
This repository includes a pre-configured GitHub Actions workflow (.github/workflows/go.yml) to automate testing and linting on every commit or pull request.

---

### License

This project is licensed under the MIT License.
