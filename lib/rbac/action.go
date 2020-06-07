package rbac

import (
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
