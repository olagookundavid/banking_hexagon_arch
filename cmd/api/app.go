package app

import (
	"sync"

	"github.com/olagookundavid/banking_hexagon/internal/jsonlog"
	"github.com/olagookundavid/banking_hexagon/internal/services"
)

type Application struct {
	Services services.Services
	Logger   *jsonlog.Logger
	Wg       sync.WaitGroup
}
