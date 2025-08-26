# Makefile for parrot
# Intelligent CLI command failure assistant

# Project configuration
PROJECT_NAME = parrot
VERSION = 1.0.1
TARGET = parrot

# Go configuration
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod

# Build flags
LDFLAGS = -ldflags="-w -s"
BUILD_FLAGS = $(LDFLAGS)

# Directories
SRCDIR = .
BUILDDIR = build
RPMDIR = rpmbuild
DISTDIR = dist

# RPM configuration
RPMVERSION = $(VERSION)
RPMRELEASE = 1
RPM_TOPDIR = $(shell pwd)/$(RPMDIR)
SPEC_FILE = $(PROJECT_NAME).spec

# Default target
.PHONY: all
all: build

# Build the Go binary
.PHONY: build
build:
	@echo "Building $(PROJECT_NAME)..."
	$(GOMOD) download
	$(GOBUILD) $(BUILD_FLAGS) -o $(TARGET) $(SRCDIR)
	@echo "✓ Build complete: $(TARGET)"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -f $(TARGET)
	rm -rf $(BUILDDIR)
	rm -rf $(RPMDIR)
	rm -rf $(DISTDIR)
	@echo "✓ Clean complete"

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Install locally
.PHONY: install
install: build
	@echo "Installing $(TARGET)..."
	@mkdir -p $(HOME)/.local/bin
	@mkdir -p $(HOME)/.local/share/$(PROJECT_NAME)
	@mkdir -p $(HOME)/.config/$(PROJECT_NAME)
	@cp $(TARGET) $(HOME)/.local/bin/
	@cp parrot-hook.sh $(HOME)/.local/share/$(PROJECT_NAME)/
	@cp config/parrot.toml.example $(HOME)/.config/$(PROJECT_NAME)/ 2>/dev/null || true
	@echo "✓ Installed to ~/.local/bin/$(TARGET)"
	@echo "Make sure ~/.local/bin is in your PATH"
	@echo "Run '$(TARGET) setup' to complete configuration"

# Uninstall locally
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(PROJECT_NAME)..."
	@rm -f $(HOME)/.local/bin/$(TARGET)
	@rm -rf $(HOME)/.local/share/$(PROJECT_NAME)
	@echo "✓ Uninstalled (config preserved in ~/.config/$(PROJECT_NAME))"

# Create source tarball for RPM
.PHONY: tarball
tarball: clean
	@echo "Creating source tarball..."
	@mkdir -p $(DISTDIR)
	@git archive --format=tar.gz --prefix=$(PROJECT_NAME)-$(VERSION)/ HEAD > $(DISTDIR)/$(PROJECT_NAME)-$(VERSION).tar.gz
	@echo "✓ Created $(DISTDIR)/$(PROJECT_NAME)-$(VERSION).tar.gz"

# Prepare RPM build environment
.PHONY: rpm-prep
rpm-prep: tarball
	@echo "Preparing RPM build environment..."
	@mkdir -p $(RPM_TOPDIR)/{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}
	@cp $(DISTDIR)/$(PROJECT_NAME)-$(VERSION).tar.gz $(RPM_TOPDIR)/SOURCES/
	@cp $(SPEC_FILE) $(RPM_TOPDIR)/SPECS/
	@echo "✓ RPM environment ready"

# Build source RPM
.PHONY: srpm
srpm: rpm-prep
	@echo "Building source RPM..."
	rpmbuild --define "_topdir $(RPM_TOPDIR)" -bs $(RPM_TOPDIR)/SPECS/$(SPEC_FILE)
	@echo "✓ Source RPM created in $(RPMDIR)/SRPMS/"

# Build binary RPM
.PHONY: rpm
rpm: rpm-prep
	@echo "Building RPM package..."
	rpmbuild --define "_topdir $(RPM_TOPDIR)" -ba $(RPM_TOPDIR)/SPECS/$(SPEC_FILE)
	@echo "✓ RPM packages created:"
	@find $(RPMDIR)/RPMS -name "*.rpm" -exec echo "  {}" \;
	@find $(RPMDIR)/SRPMS -name "*.rpm" -exec echo "  {}" \;

# Copy RPMs to repository structure (matching existing projects)
.PHONY: copy-rpms
copy-rpms: rpm
	@echo "Copying RPMs to repository structure..."
	@mkdir -p ../repos-musicsian-com/RPMS/
	@cp $(RPMDIR)/RPMS/x86_64/$(PROJECT_NAME)-*.rpm ../repos-musicsian-com/RPMS/
	@cp $(RPMDIR)/SRPMS/$(PROJECT_NAME)-*.src.rpm ../repos-musicsian-com/RPMS/
	@echo "✓ RPMs copied to ../repos-musicsian-com/RPMS/"

# Format Go code
.PHONY: format
format:
	@echo "Formatting Go code..."
	@command -v gofmt >/dev/null 2>&1 && gofmt -w . || echo "gofmt not found"
	@command -v goimports >/dev/null 2>&1 && goimports -w . || echo "goimports not found"

# Lint Go code
.PHONY: lint
lint:
	@echo "Running Go linter..."
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run || echo "golangci-lint not found"

# Run development cycle
.PHONY: dev
dev: clean format lint test build

# Quick smoke test
.PHONY: smoke-test
smoke-test: build
	@echo "Running smoke test..."
	@./$(TARGET) --version >/dev/null && echo "✓ Version command works"
	@./$(TARGET) --help >/dev/null && echo "✓ Help command works"
	@./$(TARGET) status >/dev/null 2>&1 && echo "✓ Status command works" || echo "⚠ Status command needs configuration"

# Create .repo file for YUM repository
.PHONY: repo-file
repo-file:
	@echo "Creating repository file..."
	@mkdir -p $(DISTDIR)
	@echo "[$(PROJECT_NAME)]" > $(DISTDIR)/$(PROJECT_NAME).repo
	@echo "name=Parrot - Intelligent CLI Assistant" >> $(DISTDIR)/$(PROJECT_NAME).repo
	@echo "baseurl=https://repos.musicsian.com/" >> $(DISTDIR)/$(PROJECT_NAME).repo
	@echo "enabled=1" >> $(DISTDIR)/$(PROJECT_NAME).repo
	@echo "gpgcheck=0" >> $(DISTDIR)/$(PROJECT_NAME).repo
	@echo "Created $(DISTDIR)/$(PROJECT_NAME).repo"

# Show build information
.PHONY: info
info:
	@echo "Project: $(PROJECT_NAME) v$(VERSION)"
	@echo "Target: $(TARGET)"
	@echo "Go version: $(shell $(GOCMD) version)"
	@echo "RPM build directory: $(RPM_TOPDIR)"

# Check dependencies
.PHONY: deps
deps:
	@echo "Checking dependencies..."
	@echo "Required tools:"
	@command -v $(GOCMD) >/dev/null 2>&1 && echo "  $(GOCMD) - $(shell $(GOCMD) version)" || echo "  $(GOCMD) - REQUIRED"
	@command -v make >/dev/null 2>&1 && echo "  make" || echo "  make - REQUIRED"
	@echo "RPM build tools:"
	@command -v rpmbuild >/dev/null 2>&1 && echo "  rpmbuild" || echo "  rpmbuild - for RPM building"
	@echo "Optional tools:"
	@command -v gofmt >/dev/null 2>&1 && echo "  gofmt" || echo "  gofmt - for formatting"
	@command -v golangci-lint >/dev/null 2>&1 && echo "  golangci-lint" || echo "  golangci-lint - for linting"
	@command -v goimports >/dev/null 2>&1 && echo "  goimports" || echo "  goimports - for import management"

# Help target
.PHONY: help
help:
	@echo "$(PROJECT_NAME) Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  build        Build the Go binary"
	@echo "  clean        Remove build artifacts"
	@echo "  test         Run tests"
	@echo "  install      Install locally to ~/.local/bin"
	@echo "  uninstall    Remove local installation"
	@echo "  format       Format Go code"
	@echo "  lint         Run Go linter"
	@echo "  dev          Development cycle (clean + format + lint + test + build)"
	@echo "  smoke-test   Quick functionality check"
	@echo ""
	@echo "RPM Packaging:"
	@echo "  tarball      Create source tarball"
	@echo "  rpm-prep     Prepare RPM build environment"
	@echo "  srpm         Build source RPM"
	@echo "  rpm          Build binary RPM"
	@echo "  copy-rpms    Copy RPMs to repository structure"
	@echo "  repo-file    Create .repo file for YUM"
	@echo ""
	@echo "Utilities:"
	@echo "  info         Show build information"
	@echo "  deps         Check dependencies"
	@echo "  help         Show this help"
	@echo ""
	@echo "Example workflow:"
	@echo "  make dev     # Development cycle"
	@echo "  make rpm     # Build RPM packages"
	@echo "  make copy-rpms # Deploy to repository"

.PHONY: all build clean test install uninstall tarball rpm-prep srpm rpm copy-rpms format lint dev smoke-test repo-file info deps help