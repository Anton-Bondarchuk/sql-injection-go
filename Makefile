help:
	@echo "Available commands:"
	@echo "  build      - Build and start containers, then initialize database tables"
	@echo "  start      - Start the containers"
	@echo "  stop       - Stop and remove containers"
	@echo "  help       - Show this help message"

build:
	@cd docker && docker compose up
	@cd docker && docker compose exec app_sql_injection cd migrations && go run init_table.go

start:
	@cd docker && docker compose up 

stop: 
	@cd docker && docker compose down