FROM golang:1.26.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Build con ruta absoluta
RUN go build -o /app/todo /app/cmd/main.go
# Verifica que el binario se creó
RUN test -f /app/todo && echo "✅ Binary created" || (echo "❌ ERROR: Binary not found" && exit 1)
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
# Copia el binario y el directorio web completo
COPY --from=builder /app/todo /app/todo
COPY --from=builder /app/web /app/web
EXPOSE 8080
CMD ["./todo"]