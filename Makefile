IMAGE_NAME = guilliman

CONTAINER_PORT = 8080

HOST_PORT = 8080

DB_VOLUME = ./db

ENV_VARS = -e EXCHANGE_RATE_API_KEY=$(EXCHANGE_RATE_API_KEY)

.PHONY: all build run

all: build run

dev: 
	docker build -t $(IMAGE_NAME) ./Dockerfile.dev

build:
	docker build -t $(IMAGE_NAME) .

run:
	docker run -p $(HOST_PORT):$(CONTAINER_PORT) \
		$(ENV_VARS) \
		-v $(DB_VOLUME):/app/db \
		$(IMAGE_NAME)

