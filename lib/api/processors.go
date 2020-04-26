package api

import (
//"github.com/gin-gonic/gin"
)

type jwt_authenticate_body struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	Login string
	Name  string
}
