package app

import (
	"fmt"

	"github.com/glenntam/todoken/internal/service"
)

type CoreApp struct {
	Service: *service.Service
}

type WebApp struct {
	*CoreApp
}

type ApiApp struct {
	*CoreApp
}
