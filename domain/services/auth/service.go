package auth

import (
	"database/sql"
	"github.com/vickydk/gosk/domain/dao"
	"github.com/vickydk/gosk/domain/dao/impl/mysql"
	"github.com/vickydk/gosk/domain/entity"
	"github.com/vickydk/gosk/domain/entity/model"
)

type Service interface {
	Authenticate(string, string) (*entity.Response, error)
	Me(string) (*entity.Response, error)
}

func New(db *sql.DB, j TokenGenerator, sec Securer, internal *Internal) *Auth {
	return &Auth{db: db, tg: j, sec: sec, internal: Internal{UDB: internal.UDB}}
}

// Initialize initalizes User application service with defaults
func Initialize(db *sql.DB, j TokenGenerator, sec Securer) *Auth {
	return New(db, j, sec, NewInternal(db))
}

type Auth struct {
	db       *sql.DB
	tg       TokenGenerator
	sec      Securer
	internal Internal
}

type Internal struct {
	UDB  dao.UDB
}

func NewInternal(db *sql.DB) *Internal {
	return &Internal{
		UDB:  mysql.NewUsersDaoImpl(),
	}
}

// TokenGenerator represents token generator (jwt) interface
type TokenGenerator interface {
	GenerateToken(*model.Users) (string, string, error)
}

// Securer represents security interface
type Securer interface {
	Token(string) string
}
