package mysql

import (
	"database/sql"
	"fmt"
	"github.com/vickydk/gosk/domain/entity/model"
	"github.com/vickydk/gosk/utl/log"
	"github.com/vickydk/gosk/utl/structs"
	"go.uber.org/zap"
	"strings"
)

func NewUsersDaoImpl() *UsersDaoImpl {
	return &UsersDaoImpl{}
}

type UsersDaoImpl struct{}

func (u *UsersDaoImpl) Create(db *sql.Tx, req *model.CreateUserReq) error {
	var buf strings.Builder
	buf.WriteString("INSERT INTO ")
	buf.WriteString(model.TblUsers)
	buf.WriteString(" (email, password, name, phone, merchant_code, store_code, active, role_id, group_id) VALUES ")
	buf.WriteString("(?, ?, ?, ?, ?, ?, ?, ?, ?)")
	log.Debug("sql: ", zap.String("query", buf.String()), zap.Any("parameter", req))
	stmt, err := db.Prepare(buf.String())
	defer stmt.Close()
	_, err = stmt.Exec(req.Email, req.Password, req.Name, req.Phone, req.MerchantCode, req.StoreCode, req.Active, req.RoleId, req.GroupId)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (u *UsersDaoImpl) Find(db *sql.DB, p *model.Pagination, req *model.FilterUserReq) (*[]model.Users, int, error) {
	var buf strings.Builder
	buf.WriteString("SELECT id, email, password, name, phone, merchant_code, store_code, role_id, group_id, active FROM ")
	buf.WriteString(model.TblUsers)
	buf.WriteString(" WHERE 1 = ? ")
	bind := []interface{}{1}

	if len(req.Email) > 0 {
		buf.WriteString(" AND email LIKE ?")
		bind = append(bind, fmt.Sprint("%", req.Email, "%"))
	}
	if len(req.Name) > 0 {
		buf.WriteString(" AND name LIKE ?")
		bind = append(bind, fmt.Sprint("%", req.Name, "%"))
	}
	if len(req.GroupId) > 0 {
		buf.WriteString(" AND group_id LIKE ?")
		bind = append(bind, fmt.Sprint("%", req.GroupId, "%"))
	}
	if req.RoleId > 0 {
		buf.WriteString(" AND role_id = ?")
		bind = append(bind, req.RoleId)
	}

	var total int
	var buf_count strings.Builder
	buf_count.WriteString("SELECT COUNT(*) FROM ( ")
	buf_count.WriteString(buf.String())
	buf_count.WriteString(" ) tblCount ")
	log.Debug("sql count: ", zap.String("query", buf_count.String()), zap.Any("parameter", bind))
	err := db.QueryRow(buf_count.String(), bind...).Scan(&total)
	if err != nil {
		log.Error(err)
		return nil, 0, err
	}

	buf.WriteString(" LIMIT ?")
	bind = append(bind, p.Limit)
	buf.WriteString(" OFFSET ?")
	bind = append(bind, p.Offset)

	log.Debug("sql: ", zap.String("query", buf.String()), zap.Any("parameter", bind))
	rows, err := db.Query(buf.String(), bind...)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		log.Error(err)
		return nil, 0, err
	}

	var users []model.Users
	for rows.Next() {
		var each model.Users
		structs.MergeRow(rows, &each)
		users = append(users, each)
	}

	return &users, total, nil
}

func (u *UsersDaoImpl) View(db *sql.DB, req *model.UserReq) (model.Users, error) {
	var users model.Users
	var buf strings.Builder

	buf.WriteString("SELECT ")
	buf.WriteString(model.TblUsers)
	buf.WriteString(".id, email, password, name, phone, merchant_code, store_code, role_id, permission, access_level, group_id, active, token FROM ")
	buf.WriteString(model.TblUsers)
	buf.WriteString(" LEFT JOIN ")
	buf.WriteString(model.TblRole)
	buf.WriteString(" tblR on tblR.id = role_id ")
	buf.WriteString(" WHERE 1 = ? ")
	bind := []interface{}{1}
	if len(req.Email) > 0 {
		buf.WriteString("AND email = ? ")
		bind = append(bind, req.Email)
	} else if req.Id > 0 {
		buf.WriteString("AND ")
		buf.WriteString(model.TblUsers)
		buf.WriteString(".id = ? ")
		bind = append(bind, req.Id)
	}
	log.Debug("sql: ", zap.String("query", buf.String()), zap.Any("parameter", bind))
	rows, err := db.Query(buf.String(), bind...)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()

	if err != nil {
		return users, err
	}

	for rows.Next() {
		structs.MergeRow(rows, &users)
	}

	return users, nil
}

func (u *UsersDaoImpl) Update(db *sql.Tx, userId int64, setUpdate *strings.Builder, bind *[]interface{}) error {
	var buf strings.Builder
	buf.WriteString("UPDATE ")
	buf.WriteString(model.TblUsers)
	buf.WriteString(setUpdate.String())
	buf.WriteString(" WHERE id = ?")
	*bind = append(*bind, userId)
	log.Debug("sql: ", zap.String("query", buf.String()), zap.Any("parameter", *bind))
	stmt, err := db.Prepare(buf.String())
	defer stmt.Close()
	_, err = stmt.Exec(*bind...)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (u *UsersDaoImpl) Delete(db *sql.Tx, req *model.UserReq) error {
	var buf strings.Builder
	buf.WriteString("UPDATE ")
	buf.WriteString(model.TblUsers)
	buf.WriteString(" SET ACTIVE = 0 WHERE id = ? ")
	log.Debug("sql: ", zap.String("query", buf.String()), zap.Any("parameter", req))
	stmt, err := db.Prepare(buf.String())
	defer stmt.Close()
	_, err = stmt.Exec(req.Id)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
