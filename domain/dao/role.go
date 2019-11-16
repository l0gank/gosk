package dao

import (
	"database/sql"
	"github.com/vickydk/gosk/domain/entity/model"
)

type RDB interface {
	List(*sql.DB) ([]model.Role, error)
	View(*sql.DB, *model.RoleReq) (model.Role, error)
}
