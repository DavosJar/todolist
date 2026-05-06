# Todo App Multitenant

Lista de tareas con autenticación y aislamiento por tenant.

## Deploy en Render

1. Sube el código a GitHub
2. En Render, crea un nuevo Web Service
3. Conecta tu repo
4. Render detectará el Dockerfile automáticamente
5. Agrega una base de datos PostgreSQL en Render
6. Configura variables: DATABASE_URL, JWT_SECRET, PORT
7. Deploy

## Desarrollo local

```
make start
make stop
```
