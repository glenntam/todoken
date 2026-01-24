// Package router is a custom CRUD router based on
// go-chi, scs and sqlc.
package router

import (
	"database/sql"
	"html/template"

	"github.com/glenntam/todoken/internal/models"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

type Router struct {
	chi.Router
	models    *models.Queries
	sessions  *scs.SessionManager
	templates *template.Template
}

// NewRouter return a complete new router that embeds a database,
// connection, models, session manager and template.
//
// tpls is a string glob to a directory of templates.
// conn is a DB connection.
func NewRouter(tpls string, conn *sql.DB) *Router {
	t := template.Must(template.ParseGlob(tpls))
	sm := newSessionManager()

	r := &Router{
		Router:    chi.NewRouter(),
		models:    models.New(conn),
		sessions:  sm,
		templates: t,
	}

	r.Use(r.sessions.LoadAndSave)
	r.Get("/login", r.loginForm)
	r.Post("/login", r.login)
	r.Get("/register", r.registerForm)
	r.Post("/register", r.register)
	r.Post("/logout", r.logout)

	r.Group(func(cr chi.Router) {
		cr.Use(r.requireAuth)
		cr.Get("/", r.dashboard)
	})

	return r
}
