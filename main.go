package main

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tylerb/graceful"

	"github.com/digorithm/meal_planner/application"
	"github.com/digorithm/meal_planner/finchgo"
	"github.com/digorithm/meal_planner/handlers"
	"github.com/digorithm/meal_planner/models"
)

func init() {
	gob.Register(&models.UserRow{})
}

func newConfig() (*viper.Viper, error) {
	DB_ADDR := "db"

	c := viper.New()
	c.SetDefault("dsn", fmt.Sprintf("postgres://%v:%v@%v:5432/meal_planner?sslmode=disable", "postgres", "123", DB_ADDR))
	c.SetDefault("cookie_secret", "bginXRnaDjqwiOwb")
	c.SetDefault("http_addr", ":8888")
	c.SetDefault("http_cert_file", "")
	c.SetDefault("http_key_file", "")
	c.SetDefault("http_drain_interval", "1s")

	c.AutomaticEnv()

	return c, nil
}

func Initialize() {

	// Create new instance of Finch, make it global
	Finch := finchgo.NewFinch("finch_knobs.json", "finch_sla.json")
	handlers.Finch = Finch
	application.Finch = Finch

	config, err := newConfig()
	if err != nil {
		logrus.Fatal(err)
	}

	app, err := application.New(config)
	if err != nil {
		logrus.Fatal(err)
	}

	middle, err := app.MiddlewareStruct()
	if err != nil {
		logrus.Fatal(err)
	}

	serverAddress := config.Get("http_addr").(string)

	certFile := config.Get("http_cert_file").(string)
	keyFile := config.Get("http_key_file").(string)
	drainIntervalString := config.Get("http_drain_interval").(string)

	drainInterval, err := time.ParseDuration(drainIntervalString)
	if err != nil {
		logrus.Fatal(err)
	}

	srv := &graceful.Server{
		Timeout: drainInterval,
		Server:  &http.Server{Addr: serverAddress, Handler: middle},
	}

	logrus.Infoln("Running HTTP server on " + serverAddress)

	if certFile != "" && keyFile != "" {
		err = srv.ListenAndServeTLS(certFile, keyFile)
	} else {
		err = srv.ListenAndServe()
	}

	if err != nil {
		logrus.Fatal(err)
	}

}

func main() {
	Initialize()
}
