package model

const (
	TblRole string = "role"
)

// Role model
type Role struct {
	ID          int64  `json:"id"`
	AccessLevel string `json:"access_level"`
	Desc        string `json:"description"`
	Permission  string `json:"permission"`
}

type RoleReq struct {
	Id int64 `json:"id" validate:"required"`
}
