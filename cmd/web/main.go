package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"
	//"time"

	//"github.com/glenntam/todoken/internal/database"
	//"github.com/glenntam/todoken/internal/env"
	"github.com/glenntam/todoken/internal/session"
	"github.com/glenntam/todoken/internal/smtp"
	"github.com/glenntam/todoken/internal/version"

	//"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/glenntam/envwrapper"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

type config struct {
	baseURL  string
	httpPort int
	cookie   struct {
		secretKey string
	}
	db struct {
		dsn         string
		automigrate bool
	}
	notifications struct {
		email string
	}
	session struct {
		cookieName string
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		from     string
	}
}

type application struct {
	config         config
	//db             *database.DB
	db             *sql.DB
	logger         *slog.Logger
	mailer         *smtp.Mailer
	sessionManager *scs.SessionManager
	wg             sync.WaitGroup
}

func run(logger *slog.Logger) error {
	// Environment variables:
	c := map[string]any{
		"TODOKEN_TIMEZONE": "UTC",
		"TODOKEN_PORT":     ":8080",
		"TODOKEN_DBFILE":   "todoken.db",
		"TODOKEN_LOGFILE":  "logfile.json",
		"BASE_URL": "http://localhost:3900",
		"HTTP_PORT": 3900,
		"COOKIE_SECRET_KEY": "wokutdsagetd44si7g5bgw6aonsgdfwk",
		"DB_DSN": "db.sqlite?_foreign_keys=on",
		"DB_AUTOMIGRATE": true,
		"NOTIFICATIONS_EMAIL": "",
		"SESSION_COOKIE_NAME": "session_7trykimt",
		"SMTP_HOST": "example.smtp.host",
		"SMTP_PORT": 25,
		"SMTP_USERNAME": "example_username",
		"SMTP_PASSWORD": "pa55word",
		"SMTP_FROM": "Example Name <no_reply@example.org>",
	}
	env := envwrapper.Parse(c)
	defer envwrapper.WipeSecrets(env)
	var cfg config

	cfg.baseURL = env.Str["BASE_URL"]
	cfg.httpPort = env.Int["HTTP_PORT"]
	cfg.cookie.secretKey = env.Str["COOKIE_SECRET_KEY"]
	cfg.db.dsn = env.Str["DB_DSN"]
	cfg.db.automigrate = env.Bool["DB_AUTOMIGRATE"]
	cfg.notifications.email = env.Str["NOTIFICATIONS_EMAIL"]
	cfg.session.cookieName = env.Str["SESSION_COOKIE_NAME"]
	cfg.smtp.host = env.Str["SMTP_HOST"]
	cfg.smtp.port = env.Int["SMTP_PORT"]
	cfg.smtp.username = env.Str["SMTP_USERNAME"]
	cfg.smtp.password = env.Str["SMTP_PASSWORD"]
	cfg.smtp.from = env.Str["SMTP_FROM"]

	showVersion := flag.Bool("version", false, "display version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}

	// DB:
	conn, err := sql.Open("sqlite", env.Str["TODOKEN_DBFILE"])
	if err != nil {
		panic("Couldn't open sqlite file")
	}
	defer conn.Close()
	// db, err := database.New(cfg.db.dsn)
	// if err != nil {
	//     return err
	// }
	// defer db.Close()

	// if cfg.db.automigrate {
	//     err = db.MigrateUp()
	//     if err != nil {
	//         return err
	//     }
	// }

	mailer, err := smtp.NewMailer(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.from)
	if err != nil {
		return err
	}

	sm := session.NewSessionManager(env.Str["SESSION_COOKIE_NAME"], conn)
	// sessionManager := scs.New()
	// sessionManager.Store = sqlite3store.New(db.DB.DB)
	// sessionManager.Lifetime = 7 * 24 * time.Hour
	// sessionManager.Cookie.Name = cfg.session.cookieName
	// sessionManager.Cookie.Secure = true

	app := &application{
		config:         cfg,
		db:             conn,
		logger:         logger,
		mailer:         mailer,
		sessionManager: sm,
	}

	return app.serveHTTP()
}
