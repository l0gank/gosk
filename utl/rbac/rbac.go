/*
Package gorbac provides a lightweight role-based access
control implementation in Golang.
For the purposes of this package:
	* an identity has one or more roles.
	* a role requests access to a permission.
	* a permission is given to a role.
Thus, RBAC has the following model:
	* many to many relationship between identities and roles.
	* many to many relationship between roles and permissions.
	* roles can have parent roles.
*/
package rbac

import (
	"database/sql"
	"errors"
	p "github.com/vickydk/gosk/domain/dao/impl/mysql"
	"github.com/vickydk/gosk/domain/entity/model"
	"strings"
	"sync"
)

var (
	// ErrRoleNotExist occurred if a role cann't be found
	ErrRoleNotExist = errors.New("Role does not exist")
	// ErrRoleExist occurred if a role shouldn't be found
	ErrRoleExist = errors.New("Role has already existed")
	empty        = struct{}{}
)

// AssertionFunc supplies more fine-grained permission controls.
type AssertionFunc func(*RBAC, string, Permission) bool

// RBAC object, in most cases it should be used as a singleton.
type RBAC struct {
	mutex       sync.RWMutex
	Roles       Roles
	Permissions map[string]Permission
	parents     map[string]map[string]struct{}
}

// New returns a RBAC structure.
// The default role structure will be used.
func new() *RBAC {
	return &RBAC{
		Roles:   make(Roles),
		parents: make(map[string]map[string]struct{}),
	}
}

var once sync.Once
var instance *RBAC

func GetRBAC() *RBAC {
	once.Do(func() {
		instance = new()
	})
	return instance
}

func (rbac *RBAC) LoadFIrst(db *sql.DB) error {
	rbac.Permissions = make(map[string]Permission)
	rbac.LoadPermissions(db)
	roles, _ := p.NewRoleDaoImpl().List(db)

	for _, eachRole := range roles {
		p := strings.Split(eachRole.Permission, ",")
		role := NewStdRole(eachRole.AccessLevel)
		for _, pid := range p {
			_, ok := rbac.Permissions[pid]
			if !ok {
				rbac.Permissions[pid] = NewStdPermission(pid)
			}
			role.Assign(rbac.Permissions[pid])
		}
		rbac.Add(role)
	}

	return nil
}

func (rbac *RBAC) LoadPermissions(db *sql.DB) (err error) {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	prms, _, _ := p.NewPermissionDaoImpl().Find(db, &model.Pagination{Limit: 999, Offset: 0}, &model.FilterPermissionReq{})

	for _, earchPrms := range *prms {
		_, ok := rbac.Permissions[earchPrms.PermissionCode]
		if !ok {
			rbac.Permissions[earchPrms.PermissionCode] = NewStdPermission(earchPrms.PermissionCode)
		}
	}
	rbac.Roles = Roles{}

	return
}

// SetParents bind `parents` to the role `id`.
// If the role or any of parents is not existing,
// an error will be returned.
func (rbac *RBAC) SetParents(id string, parents []string) error {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	if _, ok := rbac.Roles[id]; !ok {
		return ErrRoleNotExist
	}
	for _, parent := range parents {
		if _, ok := rbac.Roles[parent]; !ok {
			return ErrRoleNotExist
		}
	}
	if _, ok := rbac.parents[id]; !ok {
		rbac.parents[id] = make(map[string]struct{})
	}
	for _, parent := range parents {
		rbac.parents[id][parent] = empty
	}
	return nil
}

// GetParents return `parents` of the role `id`.
// If the role is not existing, an error will be returned.
// Or the role doesn't have any parents,
// a nil slice will be returned.
func (rbac *RBAC) GetParents(id string) ([]string, error) {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	if _, ok := rbac.Roles[id]; !ok {
		return nil, ErrRoleNotExist
	}
	ids, ok := rbac.parents[id]
	if !ok {
		return nil, nil
	}
	var parents []string
	for parent := range ids {
		parents = append(parents, parent)
	}
	return parents, nil
}

// SetParent bind the `parent` to the role `id`.
// If the role or the parent is not existing,
// an error will be returned.
func (rbac *RBAC) SetParent(id string, parent string) error {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	if _, ok := rbac.Roles[id]; !ok {
		return ErrRoleNotExist
	}
	if _, ok := rbac.Roles[parent]; !ok {
		return ErrRoleNotExist
	}
	if _, ok := rbac.parents[id]; !ok {
		rbac.parents[id] = make(map[string]struct{})
	}
	var empty struct{}
	rbac.parents[id][parent] = empty
	return nil
}

// RemoveParent unbind the `parent` with the role `id`.
// If the role or the parent is not existing,
// an error will be returned.
func (rbac *RBAC) RemoveParent(id string, parent string) error {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	if _, ok := rbac.Roles[id]; !ok {
		return ErrRoleNotExist
	}
	if _, ok := rbac.Roles[parent]; !ok {
		return ErrRoleNotExist
	}
	delete(rbac.parents[id], parent)
	return nil
}

// Add a role `r`.
func (rbac *RBAC) Add(r Role) (err error) {
	rbac.mutex.Lock()
	if _, ok := rbac.Roles[r.ID()]; !ok {
		rbac.Roles[r.ID()] = r
	} else {
		err = ErrRoleExist
	}
	rbac.mutex.Unlock()
	return
}

func (rbac *RBAC) AddRole(r string) (err error) {
	role := NewStdRole(r)
	rbac.Add(role)
	return
}

func (rbac *RBAC) AddPermission(p string) (err error) {
	rbac.mutex.Lock()
	defer rbac.mutex.Unlock()
	_, ok := rbac.Permissions[p]
	if !ok {
		rbac.Permissions[p] = NewStdPermission(p)
	}
	return
}

// Update perms
func (rbac *RBAC) UpdateRolePermission(role, permission, action string) (err error) {
	rbac.mutex.Lock()
	r, ok := rbac.Roles[role]
	if ok {
		_, ok := rbac.Permissions[permission]
		if action == "add" {
			if !ok {
				rbac.Permissions[permission] = NewStdPermission(permission)
			}
			r.Assign(rbac.Permissions[permission])
		} else if action == "remove" {
			r.Revoke(rbac.Permissions[permission])
		}
	}
	rbac.mutex.Unlock()
	return
}

// Remove the role by `id`.
func (rbac *RBAC) Remove(id string) (err error) {
	rbac.mutex.Lock()
	if _, ok := rbac.Roles[id]; ok {
		delete(rbac.Roles, id)
		for rid, parents := range rbac.parents {
			if rid == id {
				delete(rbac.parents, rid)
				continue
			}
			for parent := range parents {
				if parent == id {
					delete(rbac.parents[rid], id)
					break
				}
			}
		}
	} else {
		err = ErrRoleNotExist
	}
	rbac.mutex.Unlock()
	return
}

// Get the role by `id` and a slice of its parents id.
func (rbac *RBAC) Get(id string) (r Role, parents []string, err error) {
	rbac.mutex.RLock()
	var ok bool
	if r, ok = rbac.Roles[id]; ok {
		for parent := range rbac.parents[id] {
			parents = append(parents, parent)
		}
	} else {
		err = ErrRoleNotExist
	}
	rbac.mutex.RUnlock()
	return
}

// IsGranted tests if the role `id` has Permission `p` with the condition `assert`.
func (rbac *RBAC) IsGranted(id string, p Permission, assert AssertionFunc) (rslt bool) {
	rbac.mutex.RLock()
	rslt = rbac.isGranted(id, p, assert)
	rbac.mutex.RUnlock()
	return
}

func (rbac *RBAC) isGranted(id string, p Permission, assert AssertionFunc) bool {
	if assert != nil && !assert(rbac, id, p) {
		return false
	}
	return rbac.recursionCheck(id, p)
}

func (rbac *RBAC) recursionCheck(id string, p Permission) bool {
	if role, ok := rbac.Roles[id]; ok {
		if role.Permit(p) {
			return true
		}
		if parents, ok := rbac.parents[id]; ok {
			for pID := range parents {
				if _, ok := rbac.Roles[pID]; ok {
					if rbac.recursionCheck(pID, p) {
						return true
					}
				}
			}
		}
	}
	return false
}
