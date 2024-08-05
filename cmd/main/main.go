package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	app "github.com/olagookundavid/banking_hexagon/cmd/api"
	"github.com/olagookundavid/banking_hexagon/internal/data"
	"github.com/olagookundavid/banking_hexagon/internal/jsonlog"
	"github.com/olagookundavid/banking_hexagon/internal/services"
)

func main() {

	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	app := &app.Application{
		Wg: sync.WaitGroup{},

		Logger:   logger,
		Services: services.NewServices(),
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8000),
		Handler:      data.Routes(app),
		ErrorLog:     log.New(logger, "", 0),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	print("started")
	err := srv.ListenAndServe()
	print(err.Error())
}
