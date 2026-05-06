package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"todo_list/config"
	"todo_list/internal/db"
	"todo_list/internal/handlers"
	"todo_list/internal/middleware"
)

func main() {
	c := config.Load()
	if c.DatabaseURL == "" {
		log.Fatal("DATABASE_URL no configurada")
	}
	d, err := db.New(c.DatabaseURL)
	if err != nil {
		log.Fatalf("Error conectando a DB: %v", err)
	}
	defer d.Close()
	r := chi.NewRouter()
	r.Get("/static/*", func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))).ServeHTTP(w, r)
	})
	auth := middleware.AuthMiddleware(c.JWTSecret, d)
	aH, tH := handlers.NewAuthHandler(d, c.JWTSecret), handlers.NewTaskHandler(d)
	r.Get("/login", aH.LoginPage)
	r.Post("/login", aH.Login)
	r.Get("/register", aH.RegisterPage)
	r.Post("/register", aH.Register)
	r.Post("/logout", aH.Logout)
	r.With(auth).Get("/", tH.List)
	r.With(auth).Post("/tasks", tH.Create)
	r.With(auth).Put("/tasks/{id}", tH.Update)
	r.With(auth).Delete("/tasks/{id}", tH.Delete)
	log.Fatal(http.ListenAndServe(":"+c.Port, r))
}
