package matchstick

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

const (
	daysPerWeek = 7
	hoursPerDay = 24
)

func newSessionManager() *scs.SessionManager {
	sm := scs.New()
	sm.Lifetime = daysPerWeek * hoursPerDay * time.Hour
	sm.Cookie.HttpOnly = true
	sm.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Persist = true
	session.Cookie.Path = "/"
	session.Cookie.Secure = false // true if HTTPS
	return sm
}
