package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

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
	Name    string
	Preview string
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
	path := c.Param("path")[1:]
	var out []NavigateItem
	if len(path) == 0 {
		for _, v := range state.SourcesList() {
			out = append(out, NavigateItem{
				Name:    v.Id,
				Preview: "/assets/img/navigate/source.svg",
			})
		}
	} else {
		source_path := strings.SplitN(path, "/", 2)
		fmt.Printf("DEBUG: naviage source: %s\n", source_path)
		if !state.SourceExists(source_path[0]) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Source not found"})
			return
		}
		// TODO: implement actual list of source
		out = append(out, NavigateItem{
			Name:    "test_folder1",
			Preview: "/assets/img/navigate/folder.svg",
		})
		out = append(out, NavigateItem{
			Name:    "test_file2",
			Preview: "/assets/img/navigate/file.svg",
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "Get navigate data", "data": out})
}
