build-fe: ## Build frontend
	@cd frontend && yarn build 

build-be:   ## Build backend 
	@go build  -o bin/main cmd/server/main.go

run:  build-fe build-be ## Run 
	@./bin/main

template: 
	@ruby ./initscript.rb $(NAME) 

help: 
	@bash ./script/print_help.sh
