package rbac

import (
	"strings"

	"github.com/euroteltr/rbac"
)

const (
	None     = rbac.None
	Create   = rbac.Create
	Read     = rbac.Read
	Update   = rbac.Update
	Delete   = rbac.Delete
	CRUD     = rbac.CRUD
	Download = rbac.Download
	Upload   = rbac.Upload

	// For self-user modification
	UpdateSelf rbac.Action = "update_self"
)

func PermSet(r *RBAC, perm_id string, http_method string) bool {
	var actions []rbac.Action

	switch http_method {
	case "GET":
		actions = append(actions, Read)
	case "POST":
		actions = append(actions, Create, Update)

		// Additional UpdateSelf for user
		if strings.HasPrefix(perm_id, "/api/v1/user") {
			actions = append(actions, UpdateSelf)
		}
	case "DELETE":
		actions = append(actions, Delete)
	}

	if r.IsPermissionExist(perm_id, None) {
		for _, action := range actions {
			if _, ok := r.GetPermission(perm_id).Load(action); !ok {
				r.GetPermission(perm_id).Store(action, nil)
				return true
			}
		}
	} else {
		r.RegisterPermission(perm_id, "Default perm", actions...)
		return true
	}

	return false
}
