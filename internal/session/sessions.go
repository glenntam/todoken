package session

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
)

const (
	daysPerWeek = 7
	hoursPerDay = 24
)

func NewSessionManager(cookieName string, db *sql.DB) *scs.SessionManager {
	sm := scs.New()
	sm.Store = sqlite3store.New(db)
	sm.Lifetime = daysPerWeek * hoursPerDay * time.Hour
	sm.Cookie.Name = cookieName
	sm.Cookie.HttpOnly = true
	sm.Cookie.SameSite = http.SameSiteLaxMode
	sm.Cookie.Persist = true
	sm.Cookie.Path = "/"
	sm.Cookie.Secure = false // true if HTTPS
	return sm
}
