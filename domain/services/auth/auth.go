package auth

import (
	"github.com/vickydk/gosk/domain/entity"
	"github.com/vickydk/gosk/domain/entity/model"
	"github.com/vickydk/gosk/utl/config"
	"github.com/vickydk/gosk/utl/constants"
	"github.com/vickydk/gosk/utl/secure"
	"github.com/vickydk/gosk/utl/structs"
	"net/http"
	"strings"
)

func (a *Auth) Authenticate(user, pass string) (*entity.Response, error) {
	var resp *entity.Response
	u, err := a.internal.UDB.View(a.db, &model.UserReq{Email: user})
	if err != nil {
		resp = entity.Respond(err, nil, http.StatusInsufficientStorage)
		return resp, nil
	}

	raws, _ := secure.Decode([]byte(u.Password))
	ok, _ := raws.Verify([]byte(pass))
	if !ok {
		resp = entity.Respond(err, nil, http.StatusUnauthorized)
		return resp, nil
	}

	if u.Active != constants.Active {
		resp = entity.Respond(err, nil, http.StatusUnauthorized)
		return resp, nil
	}

	token, expire, err := a.tg.GenerateToken(&u)
	if err != nil {
		resp = entity.Respond(err, nil, http.StatusUnauthorized)
		return resp, nil
	}

	var setUpdate strings.Builder
	bind := []interface{}{}
	structs.DifSqlSet(&u, &model.UpdateUserReq{Id: u.Id, Token: a.sec.Token(token)}, &setUpdate, &bind)

	if len(setUpdate.String()) > 0 {
		tx, errs := a.db.Begin()
		defer func() {
			if errs != nil {
				tx.Rollback()
			}
		}()

		if err = a.internal.UDB.Update(tx, u.Id, &setUpdate, &bind); err != nil {
			tx.Rollback()
			return entity.Respond(err, nil, http.StatusInsufficientStorage), err
		}

		tx.Commit()
	}

	authRep := &model.AuthToken{Token: token, Expires: expire, RefreshToken: u.Token}

	resp = entity.Respond(err, authRep, http.StatusOK)
	return resp, nil
}

func (a *Auth) Me(email string) (*entity.Response, error) {
	var resp *entity.Response
	u, err := a.internal.UDB.View(a.db, &model.UserReq{Email: email})
	if err != nil {
		resp = entity.Respond(err, nil, http.StatusInsufficientStorage)
		return resp, nil
	}

	if u.Id == 0 {
		resp = entity.Respond(err, nil, http.StatusNotFound)
		return resp, nil
	}
	if !config.Env.Debug {
		u.Password = ""
		u.Id = 0
		u.Token = ""
		u.RoleId = 0
	}
	resp = entity.Respond(err, u, http.StatusOK)

	return resp, nil
}
