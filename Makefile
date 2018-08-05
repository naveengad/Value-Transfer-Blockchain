.PHONY: all dev clean build env-up env-down run

all: clean build env-up run

dev: build run

##### BUILD
build:
	@echo "Build ..."
	@http_proxy="http://www-proxy-adcq7-new.us.oracle.com:80" https_proxy="http://www-proxy-adcq7-new.us.oracle.com:80" dep ensure
	@go build
	@echo "Build done"

##### ENV
env-up:
	@echo "Start environment ..."
	@cd network; docker-compose up --force-recreate -d
	@echo "Sleep 15 seconds in order to let the environment setup correctly"
	@sleep 15
	@echo "Environment up"

env-down:
	@echo "Stop environment ..."
	@cd network; docker-compose down
	@echo "Environment down"

##### RUN
run:
	@echo "Start app ..."
	@./balance_transfer

##### CLEAN
clean: env-down
	@echo "Clean up ..."
	@cd network; rm -rf /tmp/balance_transfer-service-* heroes-service
	@cd network; docker rm -f -v `docker ps -a --no-trunc | grep "balance_transfer-service" | cut -d ' ' -f 1` 2>/dev/null || true
	@cd network; docker rmi `docker images --no-trunc | grep "balance_transfer-service" | cut -d ' ' -f 1` 2>/dev/null || true
	@echo "Clean up done"
