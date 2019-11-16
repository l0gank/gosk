package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/vickydk/gosk/utl/config"
	"github.com/vickydk/gosk/utl/log"
	"strings"
	"time"
)

func New() (*sql.DB, error) {
	db, err := sql.Open(config.Env.DBType, initDSN())
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(config.Env.MaxOpenConn)
	db.SetMaxIdleConns(config.Env.MaxIdle)
	err = db.Ping()
	if err != nil {
		log.Info("db is not connected")
		log.Error(err)
		return nil, err
	}
	return db, nil
}

func initDSN() string {
	var psn strings.Builder
	psn.WriteString(config.Env.DBUser)
	psn.WriteString(":")
	psn.WriteString(config.Env.DBPass)
	psn.WriteString("@tcp(")
	psn.WriteString(config.Env.DBHost)
	psn.WriteString(")/")
	psn.WriteString(config.Env.DBName)
	psn.WriteString("?parseTime=true")
	return psn.String()
}