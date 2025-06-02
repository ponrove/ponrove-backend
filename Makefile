.PHONY: all

all: build

build: docker-build-server

docker-build-server:
	@echo "Building server..."
	docker build -t github.com/ponrove/ponrove-backend:latest -f ./build/docker/Dockerfile.server .
