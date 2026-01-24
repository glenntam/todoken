package router

import (
	"context"
	"net/http"
	"time"

	"github.com/glenntam/todoken/internal/models"

	"github.com/alexedwards/argon2id"
)

func (router *Router) registerForm(w http.ResponseWriter, r *http.Request) {
	router.templates.ExecuteTemplate(w, "register.html", nil)
}

func (router *Router) register(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		http.Error(w, "server error", 500)
		return
	}

	err = router.models.CreateUser(context.TODO(), models.CreateUserParams{
		Email:        email,
		PasswordHash: hash,
		CreatedAt:    time.Now().Format(time.RFC3339),
	})
	if err != nil {
		http.Error(w, "user exists", 400)
		return
	}

	http.Redirect(w, r, "/login", http.StatusFound)
}

func (router *Router) loginForm(w http.ResponseWriter, r *http.Request) {
	router.templates.ExecuteTemplate(w, "login.html", nil)
}

func (router *Router) login(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := router.models.GetUserByEmail(context.TODO(), email)
	if err != nil {
		http.Error(w, "invalid login", 401)
		return
	}

	ok, err := argon2id.ComparePasswordAndHash(password, user.PasswordHash)
	if err != nil || !ok {
		http.Error(w, "invalid login", 401)
		return
	}

	router.sessions.Put(r.Context(), "userID", user.ID)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (router *Router) logout(w http.ResponseWriter, r *http.Request) {
	router.sessions.Remove(r.Context(), "userID")
	http.Redirect(w, r, "/login", http.StatusFound)
}

func (router *Router) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !router.sessions.Exists(r.Context(), "userID") {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (router *Router) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Authenticated dashboard"))
}
