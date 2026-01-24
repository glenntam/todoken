// Package main is a todo web app.
package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/glenntam/todoken/internal/router"

	"github.com/glenntam/envwrapper"
	"github.com/glenntam/multislog"
	_ "modernc.org/sqlite"
)

func main() {
	// Environment variables:
	cfg := map[string]any{
		"TODOKEN_TIMEZONE": "UTC",
		"TODOKEN_PORT":     "8080",
		"TODOKEN_DBFILE":   "todoken.db",
		"TODOKEN_LOGFILE":  "logfile.json",
	}
	env := envwrapper.Parse(cfg)

	// Logger:
	msl := multislog.New(
		multislog.EnableTimezone(env.Str["TODOKEN_TIMEZONE"]),
		multislog.EnableConsole(slog.LevelDebug),
		multislog.EnableLogFile(slog.LevelDebug, env.Str["TODOKEN_LOGFILE"], true, true),
	)
	defer msl.Close()
	slog.SetDefault(msl.Logger)

	// DB:
	conn, err := sql.Open("sqlite", env.Str["TODOKEN_DBFILE"])
	if err != nil {
		panic("Couldn't open sqlite file")
	}
	schema, err := os.ReadFile("./sqlc/schema.sql")
	if err != nil {
		panic("Couldn't read sqlc schema file")
	}
	_, err = conn.Exec(string(schema))
	if err != nil {
		panic("Error trying to execute schema creation")
	}

	// Router:
	r := router.NewRouter("templates/*.gohtml", conn)

	// Start Webserver:
	slog.Info("listening on :" + env.Str["TODOKEN_PORT"])
	err = http.ListenAndServe(":"+env.Str["TODOKEN_PORT"], r)
	if err != nil {
		panic("Couldn't start http server")
	}
}
