package rbac

import (
	"log"

	"github.com/euroteltr/rbac"
)

const (
	admin_role_id = "admin"
	guest_role_id = "guest"
)

type RBAC struct {
	*rbac.RBAC
}

func New() *RBAC {
	return &RBAC{rbac.New(rbac.NewConsoleLogger())}
}

func (r *RBAC) Init() {
	log.Println("[INFO] Registering default RBAC roles")

	// Creating admin role
	var admin_role *rbac.Role
	if r.IsRoleExist(admin_role_id) {
		admin_role = r.GetRole(admin_role_id)
	} else {
		admin_role, _ = r.RegisterRole(admin_role_id, "Admin default role")
	}

	// Giving all the permissions to admin role
	for _, perm := range r.Permissions() {
		r.Permit(admin_role.ID, perm, perm.Actions()...)
	}

	// Creating guest role
	if !r.IsRoleExist(guest_role_id) {
		r.RegisterRole(guest_role_id, "Guest default role")
		// TODO: use api permission, not a static string not connected to api
		r.Permit(guest_role_id, r.GetPermission("/api/v1/navigate"), Read)
	}

	// Don't need to save, because this roles are in-memory
}
