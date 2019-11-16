package users

import (
	"github.com/vickydk/gosk/domain/entity"
	"github.com/vickydk/gosk/domain/entity/model"
	"github.com/vickydk/gosk/utl/config"
	"github.com/vickydk/gosk/utl/log"
	"github.com/vickydk/gosk/utl/secure"
	"github.com/vickydk/gosk/utl/structs"
	"github.com/vickydk/gosk/utl/validation"
	"net/http"
	"strings"
)

func (u *Users) Create(req *model.CreateUserReq) (*entity.Response, error) {
	if err := validation.CheckUserManagement(req); err != nil {
		log.Error(err)
		return entity.Respond(err, nil, http.StatusBadRequest), nil
	}

	users, err := u.internal.UDB.View(u.db, &model.UserReq{Email: req.Email})
	if err != nil {
		log.Error(err)
		return entity.Respond(err, nil, http.StatusInsufficientStorage), err
	}
	if users.Id > 0 {
		return entity.Respond(err, nil, http.StatusConflict), nil
	}

	tx, errs := u.db.Begin()
	defer func() {
		if errs != nil {
			tx.Rollback()
		}
	}()

	secArgon2 := secure.DefaultConfig()
	raw, err := secArgon2.Hash([]byte(req.Password), nil)
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return entity.Respond(err, nil, http.StatusInsufficientStorage), err
	}
	req.Password = string(raw.Encode())

	err = u.internal.UDB.Create(tx, req)
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return entity.Respond(err, nil, http.StatusInsufficientStorage), err
	}

	tx.Commit()

	return entity.Respond(err, nil, http.StatusOK), nil
}

func (u *Users) Find(p *model.Pagination, req *model.FilterUserReq) ([]model.Users, int, error) {
	var userResp []model.Users

	users, total, err := u.internal.UDB.Find(u.db, p, req)
	if err != nil {
		log.Error(err)
		return userResp, 0, err
	}

	for _, eachuser := range *users {
		role, _ := u.internal.RDB.View(u.db, &model.RoleReq{Id: eachuser.RoleId})
		eachuser.AccessLevel = role.AccessLevel
		eachuser.Permission = role.Permission
		userResp = append(userResp, eachuser)
	}

	return userResp, total, nil
}

func (u *Users) View(req *model.UserReq) (*entity.Response, error) {
	var resp *entity.Response
	users, err := u.internal.UDB.View(u.db, req)
	if err != nil {
		log.Error(err)
		return entity.Respond(err, nil, http.StatusInsufficientStorage), err
	}
	if !config.Env.Debug {
		users.Password = ""
	}
	resp = entity.Respond(err, users, http.StatusOK)

	return resp, nil
}

func (u *Users) Update(req *model.UpdateUserReq) (*entity.Response, error) {
	users, err := u.internal.UDB.View(u.db, &model.UserReq{Id: req.Id})
	if err != nil {
		return entity.Respond(err, nil, http.StatusInsufficientStorage), err
	}
	if users.Id == 0 {
		return entity.Respond(err, nil, http.StatusConflict), nil
	}
	if len(req.Password) > 0 {
		if req.Password != req.PasswordConfirm {
			return entity.Respond(err, nil, http.StatusConflict), nil
		}
		if err := validation.VerifyPassword(req.Password); err != nil {
			log.Error(err)
			return entity.Respond(err, nil, http.StatusBadRequest), nil
		}
	}

	if len(req.Phone) > 0 {
		if err := validation.ValidatePhoneFormat(req.Phone); err != nil {
			log.Error(err)
			return entity.Respond(err, nil, http.StatusBadRequest), nil
		}
	}

	var setUpdate strings.Builder
	bind := []interface{}{}
	structs.DifSqlSet(&users, req, &setUpdate, &bind)

	if len(setUpdate.String()) > 0 {
		tx, errs := u.db.Begin()
		defer func() {
			if errs != nil {
				tx.Rollback()
			}
		}()

		if err = u.internal.UDB.Update(tx, users.Id, &setUpdate, &bind); err != nil {
			tx.Rollback()
			return entity.Respond(err, nil, http.StatusInsufficientStorage), err
		}

		tx.Commit()
	}

	return entity.Respond(err, nil, http.StatusOK), nil
}

func (u *Users) Delete(req *model.UserReq) (*entity.Response, error) {
	users, err := u.internal.UDB.View(u.db, req)
	if err != nil {
		return entity.Respond(err, nil, http.StatusInsufficientStorage), err
	}
	if users.Id == 0 {
		return entity.Respond(err, nil, http.StatusConflict), nil
	}
	tx, errs := u.db.Begin()
	defer func() {
		if errs != nil {
			log.Error(err)
			tx.Rollback()
		}
	}()

	err = u.internal.UDB.Delete(tx, req)
	if err != nil {
		log.Error(err)
		tx.Rollback()
		return entity.Respond(err, nil, http.StatusInsufficientStorage), err
	}
	tx.Commit()

	return entity.Respond(err, nil, http.StatusOK), nil
}
