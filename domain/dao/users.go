package dao

import (
	"database/sql"
	"github.com/vickydk/gosk/domain/entity/model"
	"strings"
)

type UDB interface {
	Create(*sql.Tx, *model.CreateUserReq) error
	Find(*sql.DB, *model.Pagination, *model.FilterUserReq) (*[]model.Users, int, error)
	View(*sql.DB, *model.UserReq) (model.Users, error)
	Update(*sql.Tx, int64, *strings.Builder, *[]interface{}) error
	Delete(*sql.Tx, *model.UserReq) error
}
