package users

import (
	"database/sql"
	"github.com/vickydk/gosk/domain/dao"
	"github.com/vickydk/gosk/domain/dao/impl/mysql"
	"github.com/vickydk/gosk/domain/entity"
	"github.com/vickydk/gosk/domain/entity/model"
)

type Service interface {
	Create(*model.CreateUserReq) (*entity.Response, error)
	Update(*model.UpdateUserReq) (*entity.Response, error)
	Find(*model.Pagination, *model.FilterUserReq) ([]model.Users, int, error)
	View(*model.UserReq) (*entity.Response, error)
	Delete(*model.UserReq) (*entity.Response, error)
}

// New creates new user application service
func New(db *sql.DB, internal *Internal) *Users {
	return &Users{db: db, internal: Internal{UDB: internal.UDB, RDB: internal.RDB}}
}

// Initialize initalizes User application service with defaults
func Initialize(db *sql.DB) *Users {
	return New(db, NewInternal())
}

type Users struct {
	db       *sql.DB
	internal Internal
}

type Internal struct {
	UDB dao.UDB
	RDB dao.RDB
}

func NewInternal() *Internal {
	return &Internal{
		UDB: mysql.NewUsersDaoImpl(),
		RDB: mysql.NewRoleDaoImpl(),
	}
}
