SERVER_BINARY_NAME=ss-server
BINARY_DIR=bin
DIST_DIR=web/dist

.PHONY: all server run-server css

server: css
	@go build -o $(BINARY_DIR)/$(SERVER_BINARY_NAME) -v ./cmd/server

run-server: server
	@./$(BINARY_DIR)/$(SERVER_BINARY_NAME) --debug true

css:
	@npx tailwindcss -i ${DIST_DIR}/main.css -o ${DIST_DIR}/tailwind.css --minify --config tailwind.config.js
