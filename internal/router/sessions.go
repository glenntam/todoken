package router

import (
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

func newSessionManager() *scs.SessionManager {
	sm := scs.New()
	sm.Lifetime = 7 * 24 * time.Hour
	sm.Cookie.HttpOnly = true
	sm.Cookie.SameSite = http.SameSiteLaxMode
	return sm
}
