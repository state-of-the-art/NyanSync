package api

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/state-of-the-art/NyanSync/lib/processors"
	"github.com/state-of-the-art/NyanSync/lib/rbac"
	"github.com/state-of-the-art/NyanSync/lib/state"
)

type User struct {
	Login       string
	Name        string
	Manager     string
	Role        string
	Password    string
	PasswordNew string
}

type NavigateItem struct {
	processors.FileSystemItem
	Preview string
}

func toAPIUser(user *state.User) (api_user User) {
	api_user.Login = user.Login
	api_user.Name = user.Name
	api_user.Manager = user.Manager
	return
}

func UserGetList(c *gin.Context) {
	query := c.Query("q")
	users := state.UsersList()
	var out_users []User
	for _, u := range users {
		if strings.Contains(u.Login, query) || strings.Contains(u.Name, query) {
			out_users = append(out_users, toAPIUser(&u))
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get users list", "data": out_users})
}

func UserGet(c *gin.Context) {
	login := c.Param("login")
	if !state.UserExists(login) {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	user := state.UserGet(login)
	c.JSON(http.StatusOK, gin.H{"message": "Get user", "data": toAPIUser(&user)})
}

func UserPost(c *gin.Context) {
	acc := ContextAccount(c)
	r := state.GetRBAC()
	perm := r.GetPermission(getPermId(c.FullPath()))

	var data User
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Wrong request body: %v", err)})
		return
	}

	if data.Login == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Login can't be empty"})
		return
	}

	user := state.UserGet(c.Param("login"))
	if state.UserExists(c.Param("login")) {
		// MODIFY USER

		if user.Login == acc.Login {
			// Check access for self update
			if !r.IsGranted(acc.Role, perm, rbac.UpdateSelf) {
				c.JSON(http.StatusForbidden, gin.H{"message": "No update access"})
				return
			}

			// Check if the password is correct to edit the user
			if !user.PasswordCheck(data.Password) {
				c.JSON(http.StatusBadRequest, gin.H{"message": "Wrong password to modify user"})
				return
			}
			// Remove init flag on self-modification
			user.Init = false
			// Set the new password if provided
			if data.PasswordNew != "" {
				user.PasswordSet(data.PasswordNew)
			}
			user.Name = data.Name

			// Rename user if login is changed
			if user.Login != data.Login {
				if state.UserExists(data.Login) {
					c.JSON(http.StatusBadRequest, gin.H{"message": "User with such login already exists"})
					return
				}
				old_login := user.Login
				user.Login = data.Login
				if err := user.IsValid(); err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Wrong user data: %v", err)})
					return
				}
				state.UserRemove(old_login)
				// TODO: Return new token for the user
				// TODO: Change Manager on all the related resources
			}
		} else if user.Manager == acc.Login {
			// Check access to update
			if !r.IsGranted(acc.Role, perm, rbac.Update) {
				c.JSON(http.StatusForbidden, gin.H{"message": "No update access"})
				return
			}

			user.Manager = data.Manager
			user.Role = data.Role
		} else {
			c.JSON(http.StatusForbidden, gin.H{"message": "You have no access to modify this user"})
			return
		}
	} else {
		// CREATE USER

		// Check permission
		if !r.IsGranted(acc.Role, perm, rbac.Create) {
			c.JSON(http.StatusForbidden, gin.H{"message": "No create access"})
			return
		}

		user := state.UserGet(c.Param("login"))
		if data.Login != user.Login {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Not the same login used in body and path"})
			return
		}

		// Set params
		user.Name = data.Name
		user.Manager = data.Manager
		user.Role = data.Role
		user.PasswordSet(data.Password)
	}

	// Save the user
	if err := user.Save(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Wrong user data: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User stored", "data": toAPIUser(&user)})
}

func UserDelete(c *gin.Context) {
	login := c.Param("login")
	if !state.UserExists(login) {
		c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	state.UserRemove(login)
	c.JSON(http.StatusOK, gin.H{"message": "User removed"})
}

func AccessGetList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get access list", "data": state.AccessList()})
}

func AccessGet(c *gin.Context) {
	id := c.Param("id")
	if state.AccessExists(id) {
		c.JSON(http.StatusOK, gin.H{"message": "Get access", "data": state.AccessGet(id)})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Access not found"})
}

func AccessPost(c *gin.Context) {
	acc := ContextAccount(c)
	r := state.GetRBAC()
	perm := r.GetPermission(getPermId(c.FullPath()))

	var data state.Access
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Wrong request body: %v", err)})
		return
	}
	if c.Param("id") == "/" {
		// Check create permission
		if !r.IsGranted(acc.Role, perm, rbac.Create) {
			c.JSON(http.StatusForbidden, gin.H{"message": "No create access"})
			return
		}
		data.NewId()
	} else if c.Param("id") != data.Id {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unable to change the access ID"})
		return
	} else if !r.IsGranted(acc.Role, perm, rbac.Update) { // Check update permission
		c.JSON(http.StatusForbidden, gin.H{"message": "No update access"})
		return
	}

	if err := data.Save(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Unable to save access: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Access stored", "data": data})
}

func AccessDelete(c *gin.Context) {
	id := c.Param("id")
	if !state.AccessExists(id) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Access not found"})
		return
	}
	state.AccessRemove(id)
	c.JSON(http.StatusOK, gin.H{"message": "Access removed"})
}

func SourceGetList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get sources list", "data": state.SourceList()})
}

func SourceGet(c *gin.Context) {
	id := c.Param("id")
	if state.SourceExists(id) {
		c.JSON(http.StatusOK, gin.H{"message": "Get source", "data": state.SourceGet(id)})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "Source not found"})
}

func SourcePost(c *gin.Context) {
	acc := ContextAccount(c)
	r := state.GetRBAC()
	perm := r.GetPermission(getPermId(c.FullPath()))

	var data state.Source
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Wrong request body: %v", err)})
		return
	}

	// Check permission
	if state.SourceExists(c.Param("id")) {
		if !r.IsGranted(acc.Role, perm, rbac.Update) {
			c.JSON(http.StatusForbidden, gin.H{"message": "No update access"})
			return
		}
	} else if !r.IsGranted(acc.Role, perm, rbac.Create) {
		c.JSON(http.StatusForbidden, gin.H{"message": "No create access"})
		return
	}

	if err := data.SaveRename(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Unable to save source: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Source stored", "data": data})
}

func SourceDelete(c *gin.Context) {
	id := c.Param("id")
	if !state.SourceExists(id) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Source not found"})
		return
	}
	state.SourceRemove(id)
	c.JSON(http.StatusOK, gin.H{"message": "Source removed"})
}

func NavigateGetList(c *gin.Context) {
	// Cut the "/" char from path
	p := c.Param("path")[1:]
	var out []NavigateItem
	if len(p) == 0 {
		for _, v := range state.SourceList() {
			out = append(out, NavigateItem{
				FileSystemItem: processors.FileSystemItem{
					Name: v.Id,
					Type: processors.Folder,
				},
				Preview: "/assets/img/navigate/source.svg",
			})
		}
	} else {
		source_path := strings.SplitN(p, "/", 2)
		fmt.Printf("DEBUG: naviage source: %s\n", source_path)
		if !state.SourceExists(source_path[0]) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Source not found"})
			return
		}

		source := state.SourceGet(source_path[0])
		uri, err := url.ParseRequestURI(source.Uri)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Wrong source URI"})
			return
		}
		if len(source_path) > 1 {
			uri.Path = path.Join(uri.Path, source_path[1])
		}

		list, err := processors.UriGetList(uri)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "Unable to list URI"})
			return
		}
		for _, item := range list {
			fmt.Printf("DEBUG: item: %s\n", item)
			out_item := NavigateItem{
				FileSystemItem: item,
			}
			if out_item.Type == processors.Folder {
				out_item.Preview = "/assets/img/navigate/folder.svg"
			} else {
				// TODO: generate preview
				out_item.Preview = "/assets/img/navigate/file.svg"
			}
			out = append(out, out_item)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Get navigate data", "data": out})
}
