APP := ./dist/server
HOST := localhost:8080

build:
	go build -v -o $(APP) .

run: build
	$(APP)

test:
	@echo "Login and getting token"
	@$(eval TOKEN = $(shell curl -s -X POST -d 'username=pieter' -d 'password=claerhout' $(HOST)/login | jq ".token"))
	@echo "Token: $(TOKEN)"
	@echo "Checking restricted call"
	@curl $(HOST)/restricted -H "Authorization: Bearer $(TOKEN)"
