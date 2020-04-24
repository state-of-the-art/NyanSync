package api

import (
	"log"
	"time"
	"github.com/gin-gonic/gin"
	jwt "github.com/appleboy/gin-jwt/v2"

	"github.com/state-of-the-art/NyanSync/lib/state"
)

func initAuthV1() {
	identity_key := "id"

	mw, err := jwt.New(&jwt.GinJWTMiddleware {
		Realm:       "nyansync jwt",
		Key:         []byte("secret key"), // TODO: replace
		MaxRefresh:  time.Hour,
		IdentityKey: identity_key,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			log.Println("[DEBUG]: PayloadFunc")
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identity_key: v.Name,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			log.Println("[DEBUG]: IdentityHandler")
			claims := jwt.ExtractClaims(c)
			return &User{
				Login: claims[identity_key].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			log.Println("[DEBUG]: Authenitcator")
			var data jwt_authenticate_body
			if c.ShouldBind(&data) != nil {
				return "", jwt.ErrMissingLoginValues
			}

			user := state.UserFind(data.Login)
			if user == nil || !user.CheckPassword(data.Password) {
				return nil, jwt.ErrFailedAuthentication
			}

			return &User{
				Login: user.Login,
				Name:  user.Name,
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			log.Println("[DEBUG]: Authorizator")
			// TODO: check the func
			if v, ok := data.(*User); ok && v.Login == "admin" {
				return true
			}

			return false
		},
	})

	if err != nil {
		panic("JWT Error:" + err.Error())
	}

	api_data.JWT = mw
}

func InitV1(router *gin.Engine) {
	initAuthV1()

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", api_data.JWT.LoginHandler)
			auth.GET("/refresh_token", api_data.JWT.RefreshHandler)
		}
		/*user := v1.Group("/user")
		{
			// TODO
		}*/
	}
}


type api_s struct {
	JWT    *jwt.GinJWTMiddleware // jwt object
}

var api_data = &api_s{}
