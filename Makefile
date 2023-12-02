SERVER_BINARY_NAME=ss-server
BINARY_DIR=bin
DIST_DIR=web/dist
TLS_DIR=tls

PLATFORM ?= linux/amd64
IMG ?= lorenzophys/secure-share:dev

.PHONY: server run-server css docker-build docker-run docker-push cert

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

cert:
	@openssl req -x509 -newkey rsa:4096 -keyout ${TLS_DIR}/server.key -out ${TLS_DIR}/server.crt -days 365 -nodes
