package model


const (
	TblUsers     string = "users"
)

type Users struct {
	Id           int64  `json:"id"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	MerchantCode string `json:"merchant_code"`
	StoreCode    string `json:"store_code"`
	RoleId       int64  `json:"role_id"`
	AccessLevel  string `json:"access_level"`
	Permission   string `json:"permission"`
	GroupId      string `json:"group_id"`
	Token        string `json:"token"`
	Active       int    `json:"active"`
}

type UserReq struct {
	Email string `json:"email,omitempty"`
	Id    int64  `json:"id,omitempty"`
}

type FilterUserReq struct {
	Email   string `json:"email,omitempty"`
	Name    string `json:"name,omitempty"`
	GroupId string `json:"group_id,omitempty"`
	RoleId  int64  `json:"role_id,omitempty"`
	PaginationReq
}

type CreateUserReq struct {
	Email           string `json:"email" validate:"required"`
	Password        string `json:"password" validate:"required"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
	Name            string `json:"name" validate:"required"`
	Phone           string `json:"phone,omitempty"`
	MerchantCode    string `json:"merchant_code,omitempty"`
	StoreCode       string `json:"store_code,omitempty"`
	RoleId          int64  `json:"role_id,omitempty"`
	GroupId         string `json:"group_id,omitempty"`
	Active          int    `json:"active" validate:"required"`
}

type UpdateUserReq struct {
	Id              int64  `json:"id" validate:"required"`
	Password        string `json:"password,omitempty"`
	PasswordConfirm string `json:"password_confirm,omitempty"`
	Name            string `json:"name,omitempty"`
	Phone           string `json:"phone,omitempty"`
	MerchantCode    string `json:"merchant_code,omitempty"`
	StoreCode       string `json:"store_code,omitempty"`
	RoleId          int    `json:"role_id,omitempty"`
	GroupId         string `json:"group_id,omitempty"`
	Active          int    `json:"active,omitempty"`
	Token           string `json:"token,omitempty"`
}
