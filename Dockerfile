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
WORKDIR /root
# Copia con ruta absoluta
COPY --from=builder /app/todo /root/todo
EXPOSE 8080
CMD ["./todo"]