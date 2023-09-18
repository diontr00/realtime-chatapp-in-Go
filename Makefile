build:   ## Build  
	go build  -o bin/main cmd/server/main.go

run:  build ## Run 
	./bin/main
	]
help: 
	@bash ./script/print_help.sh
