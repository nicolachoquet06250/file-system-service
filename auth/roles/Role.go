package roles

import (
	"filesystem_service/arrays"
	"strconv"
	"strings"
)

type RoleActions struct {
	Id   int
	Name string
}

func (r *RoleActions) GetId() int {
	return r.Id
}

func (r *RoleActions) GetName() string {
	return r.Name
}

type Role struct {
	Id               int
	Name             string
	RoleActionsIds   string
	RoleActionsNames string
	IsActive         bool
}

func (r *Role) GetId() int {
	return r.Id
}

func (r *Role) GetName() string {
	return r.Name
}

func (r *Role) GetRoleActions() []*RoleActions {
	names := strings.Split(r.RoleActionsNames, "|")
	ids := arrays.Map[string, int](strings.Split(r.RoleActionsIds, "|"), func(s string) int {
		i, _ := strconv.Atoi(s)
		return i
	})

	return arrays.Map[int, *RoleActions](arrays.Keys(ids), func(i int) *RoleActions {
		return &RoleActions{
			ids[i],
			names[i],
		}
	})
}
