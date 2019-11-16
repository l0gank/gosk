package logging

import (
	"encoding/json"
	"github.com/vickydk/gosk/domain/entity"
	"github.com/vickydk/gosk/domain/entity/model"
	"github.com/vickydk/gosk/domain/services/users"
	"github.com/vickydk/gosk/utl/log"
	"go.uber.org/zap"
	"time"
)

// New creates new user logging service
func New(svc users.Service) *LogService {
	return &LogService{
		Service: svc,
	}
}

// LogService represents user logging service
type LogService struct {
	users.Service
}

const name = "user"

func (ls *LogService) Create(req *model.CreateUserReq) (*entity.Response, error) {
	defer func(begin time.Time) {
		reqString, _ := json.Marshal(req)
		log.Info("Create request", zap.String("source", name), zap.Duration("took", time.Since(begin)), zap.String("req", string(reqString)))
	}(time.Now())
	return ls.Service.Create(req)
}

func (ls *LogService) Update(req *model.UpdateUserReq) (*entity.Response, error) {
	defer func(begin time.Time) {
		reqString, _ := json.Marshal(req)
		log.Info("Update request", zap.String("source", name), zap.Duration("took", time.Since(begin)), zap.String("req", string(reqString)))
	}(time.Now())
	return ls.Service.Update(req)
}

func (ls *LogService) Find(p *model.Pagination, req *model.FilterUserReq) ([]model.Users, int, error) {
	defer func(begin time.Time) {
		reqString, _ := json.Marshal(req)
		pString, _ := json.Marshal(p)
		log.Info("Find request", zap.String("source", name), zap.Duration("took", time.Since(begin)), zap.String("Page", string(pString)), zap.String("req", string(reqString)))
	}(time.Now())
	return ls.Service.Find(p, req)
}

func (ls *LogService) View(req *model.UserReq) (*entity.Response, error) {
	defer func(begin time.Time) {
		reqString, _ := json.Marshal(req)
		log.Info("View request", zap.String("source", name), zap.Duration("took", time.Since(begin)), zap.String("req", string(reqString)))
	}(time.Now())
	return ls.Service.View(req)
}

func (ls *LogService) Delete(req *model.UserReq) (*entity.Response, error) {
	defer func(begin time.Time) {
		reqString, _ := json.Marshal(req)
		log.Info("Delete request", zap.String("source", name), zap.Duration("took", time.Since(begin)), zap.String("req", string(reqString)))
	}(time.Now())
	return ls.Service.Delete(req)
}

