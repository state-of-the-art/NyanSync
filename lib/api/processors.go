package api

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/state-of-the-art/NyanSync/lib/processors"
	"github.com/state-of-the-art/NyanSync/lib/state"
)

type User struct {
	Login    string
	Name     string
	Manager  string
	Password string
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
	if user := state.UserFind(login); user != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Get user", "data": toAPIUser(user)})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

func UserPost(c *gin.Context) {
	var data User
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Wrong request body: %v", err)})
		return
	}

	// TODO: IMPORTANT! make sure user have access to change this user

	if data.Login == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Login can't be empty"})
		return
	}
	// TODO: edit logic with using c.Param("login")
	// TODO: change manager in all resources during renaming of the user
	user := state.UserGet(data.Login)
	// Check if the password is correct to edit the user
	if !user.PassHash.IsEmpty() && !user.CheckPassword(data.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Wrong password to modify user"})
		return
	}
	// Remove init flag when user saved twice
	if !user.PassHash.IsEmpty() {
		user.Init = false
	}
	if err := user.Set(data.Password, data.Name, data.Manager, user.Init); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Wrong user data: %v", err)})
		if user.PassHash.IsEmpty() {
			// Remove the new invalid user
			user.Remove()
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User stored", "data": toAPIUser(user)})
}

func UserDelete(c *gin.Context) {
	login := c.Param("login")
	if u := state.UserFind(login); u != nil {
		u.Remove()
		c.JSON(http.StatusOK, gin.H{"message": "User removed"})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
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
	var data state.Access
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Wrong request body: %v", err)})
		return
	}
	if c.Param("id") == "/" {
		// Create new access
		data.NewId()
	} else if c.Param("id") != data.Id {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unable to change the access ID"})
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
	var data state.Source
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Wrong request body: %v", err)})
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
