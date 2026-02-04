// Package main is a todo web app.
package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"time"

	"github.com/glenntam/todoken/internal/matchstick"
	"github.com/glenntam/todoken/internal/service"

	"github.com/glenntam/envwrapper"
	"github.com/glenntam/multislog"
	_ "modernc.org/sqlite"
)

const (
	dbStartupTimeout    = 5 * time.Second
	sqlcQueriesLocation = "./sqlc/schema.sql"
	templatesGlob       = "templates/*.gohtml"
	staticDir           = "static"
)

func main() {
	// Environment variables:
	cfg := map[string]any{
		"TODOKEN_TIMEZONE": "UTC",
		"TODOKEN_PORT":     ":8080",
		"TODOKEN_DBFILE":   "todoken.db",
		"TODOKEN_LOGFILE":  "logfile.json",
	}
	env := envwrapper.Parse(cfg)
	defer envwrapper.WipeSecrets(env)

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
	//    create tables if they don't exist
	schema, err := os.ReadFile(sqlcQueriesLocation)
	if err != nil {
		panic("Couldn't read sqlc schema file")
	}
	ctx, cancel := context.WithTimeout(context.Background(), dbStartupTimeout)
	defer cancel()
	_, err = conn.ExecContext(ctx, string(schema))
	if err != nil {
		panic("Error trying to execute schema creation")
	}

	// Services:
	s :=
	// Router:
	r := matchstick.NewRouter(templatesGlob, staticDir, s, conn)

	// Webserver:
	ws := matchstick.NewWebServer(env.Str["TODOKEN_PORT"], r)
	ws.Start()
}
