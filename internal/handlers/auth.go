package handlers

import (
	"net/http"
	"time"
	"todo_list/internal/db"
	"todo_list/internal/middleware"
	"todo_list/web/templates"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db        *db.Database
	jwtSecret string
}

func NewAuthHandler(database *db.Database, jwtSecret string) *AuthHandler {
	return &AuthHandler{db: database, jwtSecret: jwtSecret}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	templates.LoginPage("").Render(r.Context(), w)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.Header().Set("Content-Type", "text/html")
		templates.LoginPage("Error al procesar formulario").Render(r.Context(), w)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		w.Header().Set("Content-Type", "text/html")
		templates.LoginPage("Email y contraseña requeridos").Render(r.Context(), w)
		return
	}

	user, err := h.db.GetUserByEmail(r.Context(), email)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		templates.LoginPage("Error en el servidor").Render(r.Context(), w)
		return
	}

	if user == nil {
		w.Header().Set("Content-Type", "text/html")
		templates.LoginPage("Credenciales inválidas").Render(r.Context(), w)
		return
	}

	passwordHash, err := h.db.GetUserPassword(r.Context(), user.ID)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		templates.LoginPage("Error en el servidor").Render(r.Context(), w)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		w.Header().Set("Content-Type", "text/html")
		templates.LoginPage("Credenciales inválidas").Render(r.Context(), w)
		return
	}

	tenant, err := h.db.GetTenantByUserID(r.Context(), user.ID)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		templates.LoginPage("Error en el servidor").Render(r.Context(), w)
		return
	}

	if tenant == nil {
		w.Header().Set("Content-Type", "text/html")
		templates.LoginPage("Tenant no encontrado").Render(r.Context(), w)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, middleware.Claims{
		UserID:   user.ID,
		TenantID: tenant.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		templates.LoginPage("Error al generar token").Render(r.Context(), w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   86400,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) RegisterPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	templates.RegisterPage("").Render(r.Context(), w)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.Header().Set("Content-Type", "text/html")
		templates.RegisterPage("Error al procesar formulario").Render(r.Context(), w)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	if email == "" || password == "" {
		w.Header().Set("Content-Type", "text/html")
		templates.RegisterPage("Email y contraseña requeridos").Render(r.Context(), w)
		return
	}

	if password != confirmPassword {
		w.Header().Set("Content-Type", "text/html")
		templates.RegisterPage("Las contraseñas no coinciden").Render(r.Context(), w)
		return
	}

	if len(password) < 6 {
		w.Header().Set("Content-Type", "text/html")
		templates.RegisterPage("La contraseña debe tener al menos 6 caracteres").Render(r.Context(), w)
		return
	}

	existingUser, err := h.db.GetUserByEmail(r.Context(), email)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		templates.RegisterPage("Error en el servidor").Render(r.Context(), w)
		return
	}

	if existingUser != nil {
		w.Header().Set("Content-Type", "text/html")
		templates.RegisterPage("Email ya registrado").Render(r.Context(), w)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		templates.RegisterPage("Error al encriptar contraseña").Render(r.Context(), w)
		return
	}

	user, err := h.db.CreateUser(r.Context(), email, string(passwordHash))
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		templates.RegisterPage("Error al crear usuario").Render(r.Context(), w)
		return
	}

	tenantName := email
	_, err = h.db.CreateTenant(r.Context(), user.ID, tenantName)
	if err != nil {
		w.Header().Set("Content-Type", "text/html")
		templates.RegisterPage("Error al crear tenant").Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	templates.RegisterSuccessPage().Render(r.Context(), w)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
