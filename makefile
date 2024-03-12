.PHONY: install
install: ## install deps
	go mod download

.PHONY: test
test: ## run tests
	go clean --testcache
	go test ./...

.PHONY: run-server
run-server: ## Run localy
	go run cmd/server/main.go

.PHONY: run-client
run-client: ## Run localy
	go run cmd/client/main.go

.PHONY: run 
run: ## Start docker compose
	docker-compose up --abort-on-container-exit --force-recreate --build server --build client

.PHONY: proto
proto: ## Build protos
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --proto_path=. internal/actor/actor.proto && protoc -I=internal/actor --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --proto_path=. internal/remote/remote.proto
	
.PHONY: help
help: ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

