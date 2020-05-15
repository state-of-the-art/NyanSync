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

type jwt_authenticate_body struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	Login string
	Name  string
}

type NavigateItem struct {
	processors.FileSystemItem
	Preview string
}

func UsersGetList(c *gin.Context) {
	users := state.UsersList()
	var out_users []User
	for _, u := range users {
		out_users = append(out_users, User{
			Login: u.Login,
			Name:  u.Name,
		})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Get users list", "data": out_users})
}

func UserPost(c *gin.Context) {
	var data state.User
	if err := c.ShouldBind(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Wrong request body: %v", err)})
		return
	}
	/*if err := data.SaveRename(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Unable to save user: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User stored"})*/
	c.JSON(http.StatusBadRequest, gin.H{"message": "Save user not implemented"})
}

func UserDelete(c *gin.Context) {
	login := c.Param("Login")
	if u := state.UserFind(login); u != nil {
		u.Remove()
		c.JSON(http.StatusOK, gin.H{"message": "User removed"})
		return
	}
	c.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
}

func SourcesGetList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get sources list", "data": state.SourcesList()})
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

	c.JSON(http.StatusOK, gin.H{"message": "Source stored"})
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
		for _, v := range state.SourcesList() {
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
