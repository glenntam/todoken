package main

import (
	"net/http"

	"github.com/glenntam/todoken/assets"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.NotFound(app.notFound)

	mux.Use(app.recoverPanic)
	mux.Use(app.securityHeaders)
	mux.Use(app.sessionManager.LoadAndSave)
	mux.Use(app.preventCSRF)
	mux.Use(app.authenticate)

	fileServer := http.FileServer(http.FS(assets.EmbeddedFiles))
	mux.Handle("/static/*", fileServer)

	mux.Get("/", app.home)

	mux.Group(func(mux chi.Router) {
		mux.Use(app.requireAnonymousUser)

		mux.Get("/signup", app.signup)
		mux.Post("/signup", app.signup)
		mux.Get("/login", app.login)
		mux.Post("/login", app.login)
		mux.Get("/forgotten-password", app.forgottenPassword)
		mux.Post("/forgotten-password", app.forgottenPassword)
		mux.Get("/forgotten-password-confirmation", app.forgottenPasswordConfirmation)
		mux.Get("/password-reset/{plaintextToken}", app.passwordReset)
		mux.Post("/password-reset/{plaintextToken}", app.passwordReset)
		mux.Get("/password-reset-confirmation", app.passwordResetConfirmation)
	})

	mux.Group(func(mux chi.Router) {
		mux.Use(app.requireAuthenticatedUser)

		mux.Get("/restricted", app.restricted)
		mux.Post("/logout", app.logout)
	})

	return mux
}
