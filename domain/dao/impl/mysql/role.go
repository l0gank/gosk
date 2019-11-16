package mysql

import (
	"database/sql"
	"github.com/vickydk/gosk/domain/entity/model"
	"github.com/vickydk/gosk/utl/log"
	"github.com/vickydk/gosk/utl/structs"
	"go.uber.org/zap"
	"strings"
)

func NewRoleDaoImpl() *RoleDaoImpl {
	return &RoleDaoImpl{}
}

type RoleDaoImpl struct{}

func (r *RoleDaoImpl) List(db *sql.DB) ([]model.Role, error) {
	var roles []model.Role
	var buf strings.Builder
	buf.WriteString("SELECT access_level, permission FROM ")
	buf.WriteString(model.TblRole)
	buf.WriteString(" WHERE 1=1")
	rows, err := db.Query(buf.String())
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var each model.Role
		structs.MergeRow(rows, &each)
		roles = append(roles, each)
	}

	return roles, nil
}

func (u *RoleDaoImpl) View(db *sql.DB, req *model.RoleReq) (model.Role, error) {
	var userrole model.Role
	var buf strings.Builder

	buf.WriteString("SELECT id, access_level, description, permission FROM ")
	buf.WriteString(model.TblRole)
	buf.WriteString(" WHERE 1 = ? ")
	bind := []interface{}{1}
	buf.WriteString("AND id = ? ")
	bind = append(bind, req.Id)
	log.Debug("sql: ", zap.String("query", buf.String()), zap.Any("parameter", req))
	rows, err := db.Query(buf.String(), bind...)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		log.Error(err)
		return userrole, err
	}

	for rows.Next() {
		structs.MergeRow(rows, &userrole)
	}

	return userrole, nil
}
