package api

import (
	"crypto/sha1"
	"github.com/vickydk/gosk/domain/services/auth"
	"github.com/vickydk/gosk/domain/services/users"
	al "github.com/vickydk/gosk/interface/api/auth/logging"
	at "github.com/vickydk/gosk/interface/api/auth/transport"
	ul "github.com/vickydk/gosk/interface/api/v1/users/logging"
	ut "github.com/vickydk/gosk/interface/api/v1/users/transport"
	"github.com/vickydk/gosk/utl/config"
	dbhandler "github.com/vickydk/gosk/utl/dbhandler/mysql"
	"github.com/vickydk/gosk/utl/middleware/jwt"
	rbc "github.com/vickydk/gosk/utl/middleware/rbac"
	"github.com/vickydk/gosk/utl/rbac"
	"github.com/vickydk/gosk/utl/secure"
	"github.com/vickydk/gosk/utl/server/api"
)

// Start starts the API service
func Start() error {
	db, err := dbhandler.New()
	if err != nil {
		return err
	}

	rbac.GetRBAC().LoadFIrst(db)
	rbc := rbc.New()
	sec := secure.New(sha1.New())
	jwt := jwt.New(config.Env.Secret, config.Env.SigningAlgorithm, config.Env.Duration)

	e := server.New()

	at.NewHTTP(al.New(auth.Initialize(db, jwt, sec)), e, jwt.MWFunc())

	v1 := e.Group("/v1")
	v1.Use(jwt.MWFunc(), rbc.MWFunc())

	ut.NewHTTP(ul.New(users.Initialize(db)), v1)

	server.Start(e, &server.Config{
		Port:                config.Env.Port,
		ReadTimeoutSeconds:  config.Env.ReadTimeout,
		WriteTimeoutSeconds: config.Env.WriteTimeout,
		Debug:               config.Env.Debug,
	})

	return nil
}