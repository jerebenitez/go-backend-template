# Variables
DOCKER_IMAGE_NAME := server
DOCKER_TAG := latest
VENV_DIR := .venv
PYTHON := $(VENV_DIR)/bin/python
PIP := $(VENV_DIR)/bin/pip
GO_SOURCES := $(shell find . -name "*.go" -not -path "./vendor/*")
TEST_DIR := ./tests

# Markers
DOCKER_MARKER := .docker-built
VENV_MARKER := $(VENV_DIR)/.venv-created
DEPS_MARKER := $(VENV_DIR)/.deps-installed

.PHONY: test
test: $(DEPS_MARKER)
	@echo "Running pytest in tests directory..."
	cd tests && ../$(PYTHON) -m pytest

$(DOCKER_MARKER): $(GO_SOURCES) Dockerfile
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) .
	@touch $(DOCKER_MARKER)

$(VENV_MARKER):
	@echo "Creating virtual environment..."
	python3 -m venv $(VENV_DIR)
	@touch $(VENV_MARKER)

$(DEPS_MARKER): $(TEST_DIR)/requirements.txt $(VENV_MARKER) $(DOCKER_MARKER)
	@echo "Installing Python dependencies..."
	$(PIP) install --upgrade pip
	$(PIP) install -r ./tests/requirements.txt
	@touch $(DEPS_MARKER)

.PHONY: clean
clean:
	@echo "Cleaning up..."
	rm -rf $(VENV_DIR)
	rm -f $(DOCKER_MARKER)
	docker rmi $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) 2>/dev/null || true

.PHONY: rebuild
rebuild: clean test

.PHONY: docker
docker: $(DOCKER_MARKER)

.PHONY: venv
venv: $(DEPS_MARKER)

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  test     - Run integration tests (default)"
	@echo "  docker   - Build Docker image only"
	@echo "  venv     - Set up Python environment only"
	@echo "  clean    - Remove all generated files and Docker image"
	@echo "  rebuild  - Clean and run tests from scratch"
	@echo "  help     - Show this help message"
