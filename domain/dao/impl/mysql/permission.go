package mysql

import (
	"database/sql"
	"fmt"
	"github.com/vickydk/gosk/domain/entity"
	"github.com/vickydk/gosk/domain/entity/model"
	"github.com/vickydk/gosk/utl/structs"
	"strings"
)

func NewPermissionDaoImpl() *PermissionDaoImpl {
	return &PermissionDaoImpl{}
}

type PermissionDaoImpl struct{}

func (u *PermissionDaoImpl) Find(db *sql.DB, p *model.Pagination, req *model.FilterPermissionReq) (*[]model.Permission, int, error) {
	var buf strings.Builder
	buf.WriteString("SELECT id, permission_code, name, group_permission, description FROM ")
	buf.WriteString(model.TblUserPermission)
	buf.WriteString(" WHERE 1 = ? ")
	bind := []interface{}{1}

	if len(req.PermissionCode) > 0 {
		buf.WriteString(" AND permission_code LIKE ?")
		bind = append(bind, fmt.Sprint("%", req.PermissionCode, "%"))
	}
	if len(req.Name) > 0 {
		buf.WriteString(" AND name LIKE ?")
		bind = append(bind, fmt.Sprint("%", req.Name, "%"))
	}

	var total int
	var buf_count strings.Builder
	buf_count.WriteString("SELECT COUNT(*) FROM ( ")
	buf_count.WriteString(buf.String())
	buf_count.WriteString(" ) tblCount ")
	err := db.QueryRow(buf_count.String(), bind...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return nil, 0, entity.ErrNotFound
	}
	buf.WriteString(" LIMIT ?")
	bind = append(bind, p.Limit)
	buf.WriteString(" OFFSET ?")
	bind = append(bind, p.Offset)

	rows, err := db.Query(buf.String(), bind...)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		return nil, 0, err
	}

	var permission []model.Permission
	for rows.Next() {
		var each model.Permission
		structs.MergeRow(rows, &each)
		permission = append(permission, each)
	}

	return &permission, total, nil
}
