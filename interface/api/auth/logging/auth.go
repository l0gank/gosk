package logging

import (
	"encoding/json"
	"github.com/vickydk/gosk/domain/entity"
	"github.com/vickydk/gosk/domain/services/auth"
	"github.com/vickydk/gosk/utl/log"
	"go.uber.org/zap"
	"time"
)

// New creates new auth logging service
func New(svc auth.Service) *LogService {
	return &LogService{
		Service: svc,
	}
}

// LogService represents auth logging service
type LogService struct {
	auth.Service
}

const name = "auth"

// Authenticate logging
func (ls *LogService) Authenticate(user, password string) (resp *entity.Response, err error) {
	defer func(begin time.Time) {
		respString, _ := json.Marshal(resp)
		log.Info("Authenticate request", zap.String("source", name), zap.Duration("took", time.Since(begin)), zap.String("req", user), zap.String("resp", string(respString)))
	}(time.Now())
	return ls.Service.Authenticate(user, password)
}

// Me logging
func (ls *LogService) Me(email string) (resp *entity.Response, err error) {
	defer func(begin time.Time) {
		respString, _ := json.Marshal(resp)
		log.Info("Me Request", zap.String("source", name), zap.Duration("took", time.Since(begin)), zap.String("req", email), zap.String("resp", string(respString)))
	}(time.Now())
	return ls.Service.Me(email)
}
