package model

const (
	TblUserPermission string = "permission"
)

type Permission struct {
	Id              int64  `json:"id"`
	PermissionCode  string `json:"permission_code"`
	Name            string `json:"name"`
	GroupPermission string `json:"group_permission"`
	Description     string `json:"description"`
}

type PermissionReq struct {
	PermissionCode string `json:"permission_code"`
	Name           string `json:"name"`
}

type FilterPermissionReq struct {
	PermissionCode string `json:"permission_code"`
	Name           string `json:"name"`
	PaginationReq
}

type UpdatePermissionReq struct {
	Id             int64  `json:"id"  validate:"required"`
	PermissionCode string `json:"permission_code"`
	Name           string `json:"name"`
}
