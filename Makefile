.PHONY: start stop build

start:
	@echo "🐳 Iniciando Docker..."
	sudo docker-compose up -d
	@sleep 5
	@echo "📊 Creando base de datos..."
	@PGPASSWORD=postgres psql -h localhost -U postgres -c "CREATE DATABASE todos_db;" 2>/dev/null || true
	@echo "📊 Creando tablas..."
	@PGPASSWORD=postgres psql -h localhost -U postgres -d todos_db -f schema.sql 2>/dev/null || true
	@echo "✅ Base de datos lista"
	@echo ""
	@echo "🚀 Iniciando aplicación..."
	@echo "📍 Accede a: http://localhost:8080"
	@DATABASE_URL=postgres://postgres:postgres@localhost:5432/todos_db go run ./cmd/main.go

stop:
	@pkill -f "go run ./cmd/main.go" || true
	@sleep 1
	sudo docker-compose down
	@echo "✅ Detenido"

build:
	@echo "🔨 Compilando..."
	go build -o ./bin/todo ./cmd/main.go
	@echo "✅ Listo: ./bin/todo"

help:
	@echo "make start   - Levanta DB + App"
	@echo "make stop    - Detiene TODO"
	@echo "make build   - Compila binario"

.DEFAULT_GOAL := help
