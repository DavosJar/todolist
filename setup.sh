#!/bin/bash

# Script de setup para la aplicación de Lista de Tareas Multitenant
# Uso: ./setup.sh

set -e

echo "🚀 Setup - Lista de Tareas Multitenant"
echo "======================================"

# Verificar dependencias
echo "✓ Verificando dependencias..."
which go > /dev/null || { echo "Error: Go no instalado"; exit 1; }
which psql > /dev/null || { echo "Error: PostgreSQL no instalado"; exit 1; }

# Descargar dependencias
echo "✓ Descargando dependencias..."
go mod tidy

# Instalar templ
echo "✓ Instalando templ..."
go install github.com/a-h/templ/cmd/templ@latest

# Generar código de templates
echo "✓ Generando código de templates..."
templ generate ./web/templates

# Crear base de datos
echo "✓ Creando base de datos (requerirá contraseña de PostgreSQL)..."
psql -U postgres -c "CREATE DATABASE IF NOT EXISTS todos_db;"
psql -U postgres -d todos_db -f schema.sql

# Compilar
echo "✓ Compilando aplicación..."
mkdir -p ./bin
go build -o ./bin/todo ./cmd/main.go

echo ""
echo "✅ Setup completado!"
echo ""
echo "Para ejecutar la aplicación:"
echo ""
echo "  export DATABASE_URL=\"postgres://postgres@localhost/todos_db\""
echo "  ./bin/todo"
echo ""
echo "O usar:"
echo "  go run ./cmd/main.go"
echo ""
