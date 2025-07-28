package handlers

import (
	"html/template"
	"net/http"

	"calespero/internal/core/domain"
	"calespero/internal/core/ports"
)

type UserHandler struct {
	templates *template.Template
	userSvc   ports.UserService
}

func NewUserHandler(templates *template.Template, userSvc ports.UserService) *UserHandler {
	return &UserHandler{
		templates: templates,
		userSvc:   userSvc,
	}
}

func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.templates.ExecuteTemplate(w, "login.html", nil)
		return
	}

	if r.Method == "POST" {
		email := r.FormValue("email")
		password := r.FormValue("password")

		token, err := h.userSvc.AuthenticateUser(r.Context(), email, password)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Set JWT token as cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			HttpOnly: true,
		})

		http.Redirect(w, r, "/start", http.StatusSeeOther)
		return
	}
}

func (h *UserHandler) HandleNewUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.templates.ExecuteTemplate(w, "new_user.html", nil)
		return
	}

	if r.Method == "POST" {
		user := &domain.User{
			Email:    r.FormValue("email"),
			Password: r.FormValue("password"),
			Phone:    r.FormValue("phone"),
		}

		err := h.userSvc.CreateUser(r.Context(), user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
}

func (h *UserHandler) HandleStart(w http.ResponseWriter, r *http.Request) {
	// Get JWT token from cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Validate token
	_, err = h.userSvc.ValidateToken(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	h.templates.ExecuteTemplate(w, "start.html", nil)
}
