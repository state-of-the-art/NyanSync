package rbac

type Permission struct {
	Id          string
	Description string
	Actions     []string
}

type Role struct {
	Id          string
	Description string
	Permissions []Permission
}

func RoleExists(id string) bool {
	return RBAC.IsRoleExist(id)
}

func RoleGet(id string) Role {
	var role Role
	if RBAC.IsRoleExist(id) {
		r := RBAC.GetRole(id)
		role = Role{r.ID, r.Description, []Permission{}}
	} else {
		role = Role{id, "", []Permission{}}
	}

	return role
}

func (r *Role) Save() error {
	// TODO
	return nil
}

func RoleRemove(id string) {
	if RBAC.IsRoleExist(id) {
		RBAC.RemoveRole(id)
	}
	Save()
}

func (r *Role) IsValid() error {
	// TODO
	return nil
}

func RoleList() []Role {
	roles := RBAC.Roles()
	out := make([]Role, 0, len(roles))

	for _, r := range roles {
		out = append(out, Role{r.ID, r.Description, []Permission{}})
	}
	return out
}
