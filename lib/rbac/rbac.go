package rbac

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/euroteltr/rbac"

	"github.com/state-of-the-art/NyanSync/lib/config"
	"github.com/state-of-the-art/NyanSync/lib/location"
)

const (
	admin_role_id = "admin"
	guest_role_id = "guest"
)

var RBAC = rbac.New(rbac.NewConsoleLogger())

func Init() {
	log.Println("[INFO] Registering default RBAC roles")

	// Creating admin role
	var admin_role *rbac.Role
	if RBAC.IsRoleExist(admin_role_id) {
		admin_role = RBAC.GetRole(admin_role_id)
	} else {
		admin_role, _ = RBAC.RegisterRole(admin_role_id, "Admin default role")
	}

	// Giving all the permissions to admin role
	for _, perm := range RBAC.Permissions() {
		RBAC.Permit(admin_role.ID, perm, perm.Actions()...)
	}

	// Creating guest role
	if !RBAC.IsRoleExist(guest_role_id) {
		RBAC.RegisterRole(guest_role_id, "Guest default role")
		// TODO: use api permission, not a static string not connected to api
		RBAC.Permit(guest_role_id, RBAC.GetPermission("/api/v1/navigate"), Read)
	}

	// Don't need to save, because this roles are in-memory
}

func Load() {
	path := location.RealFilePath(config.Cfg().RBACFilePath)
	log.Println("[DEBUG] Loading RBAC from file", path)
	if file, err := os.OpenFile(path, os.O_RDONLY, 640); err == nil {
		defer file.Close()
		if err = RBAC.LoadJSON(file); err != nil {
			log.Panic("Unable to load rbac file: ", err)
		}
	}
}

func Save() {
	path := location.RealFilePath(config.Cfg().RBACFilePath)
	log.Println("[DEBUG] Saving RBAC to file", path)
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		log.Panic("Unable to create save dir: ", err)
	}
	if file, err := os.OpenFile(path, os.O_WRONLY, 640); err == nil {
		defer file.Close()
		if err = RBAC.SaveJSON(file); err != nil {
			log.Panic("Unable to load save file: ", err)
		}
	}
}

func PermSet(perm_id string, http_method string) bool {
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

	if RBAC.IsPermissionExist(perm_id, None) {
		for _, action := range actions {
			if _, ok := RBAC.GetPermission(perm_id).Load(action); !ok {
				RBAC.GetPermission(perm_id).Store(action, nil)
				return true
			}
		}
	} else {
		RBAC.RegisterPermission(perm_id, "Default perm", actions...)
		return true
	}

	return false
}
