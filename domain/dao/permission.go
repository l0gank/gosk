package dao

import (
	"database/sql"
	"github.com/vickydk/gosk/domain/entity/model"
)

type PDB interface {
	Find(*sql.DB, *model.Pagination, *model.FilterPermissionReq) (*[]model.Permission, int, error)
}
