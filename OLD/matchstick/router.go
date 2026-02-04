// Package matchstick is a web app boilerplate library
// based on common patterns using go-chi, scs and sqlc.
package matchstick

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/glenntam/todoken/internal/model"
	// "github.com/glenntam/todoken/internal/service"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
)

// Router contains a mux and the interdependant db models, session manager, and templates.
type Router struct {
	chi.Router

	models    *model.Queries
	sessions  *scs.SessionManager
	templates *template.Template
}

// NewRouter return a complete new router that embeds a database,
// connection, models, session manager and template.
//
// tpls is a string glob to a directory of templates.
// conn is a DB connection.
func NewRouter(tpls, staticDir string, conn *sql.DB) *Router {
	r := &Router{
		Router:    chi.NewRouter(),
		model:     model.New(conn),
		sessions:  newSessionManager(),
		templates: template.Must(template.ParseGlob(tpls)),
	}

	r.Handle(
		"/"+staticDir+"/*",
		http.StripPrefix("/"+staticDir+"/", http.FileServer(http.Dir("./"+staticDir))),
	)

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
