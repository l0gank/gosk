package dao

import (
	"database/sql"
	omsApi "omsApi/pkg/utl/model"
)

type PDB interface {
	Find(*sql.DB, *omsApi.Pagination, *omsApi.FilterPermissionReq) (*[]omsApi.Permission, int, error)
}
