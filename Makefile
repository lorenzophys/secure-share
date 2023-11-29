SERVER_BINARY_NAME=ss-server
BINARY_DIR=bin
DIST_DIR=web/dist

PLATFORM ?= linux/amd64
IMG ?= lorenzophys/secure-share:dev
PORT ?= 8080

.PHONY: all server run-server css docker-build docker-run

server: css
	@go build -o $(BINARY_DIR)/$(SERVER_BINARY_NAME) -v ./cmd/server

run-server: server
	@./$(BINARY_DIR)/$(SERVER_BINARY_NAME) --debug true

css:
	@npx tailwindcss -i ${DIST_DIR}/main.css -o ${DIST_DIR}/tailwind.css --minify --config tailwind.config.js

test:
	@ginkgo run -r -v

docker-build:
	@docker build -t ${IMG} --platform ${PLATFORM} .

docker-push:
	@docker push ${IMG}

docker-run:
	@docker run --rm -p 8080:${PORT} ${IMG}
