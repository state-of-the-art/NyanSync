package api

import (
	"net/http"

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

func GetSourcesList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Get sources list", "data": state.SourcesList()})
}
