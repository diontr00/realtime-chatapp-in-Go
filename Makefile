make-cert: ## Create Certificate
	@ruby ./script/genCert.rb

build-fe: ## Build frontend
	@cd frontend && yarn build 

build-be:   ## Build backend 
	@go build  -o bin/main cmd/server/main.go

run:  make-cert build-fe build-be ## Run 
	@./bin/main

template: 
	@ruby ./initscript.rb $(NAME) 

help: 
	@bash ./script/print_help.sh
