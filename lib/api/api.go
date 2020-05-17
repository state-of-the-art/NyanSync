package api

import (
	"log"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/state-of-the-art/NyanSync/lib/state"
)

func initAuthV1() {
	identity_key := "id"

	mw, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "nyanshare jwt",
		Key:         []byte("secret key"), // TODO: replace
		MaxRefresh:  time.Hour,
		IdentityKey: identity_key,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			log.Println("[DEBUG]: PayloadFunc")
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identity_key: v.Login,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			log.Println("[DEBUG]: IdentityHandler")
			claims := jwt.ExtractClaims(c)
			log.Println("[DEBUG]: IdentityHandler", claims)
			return &User{
				Login: claims[identity_key].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			log.Println("[DEBUG]: Authenticator")
			var data jwt_authenticate_body
			if c.ShouldBind(&data) != nil {
				return "", jwt.ErrMissingLoginValues
			}

			user := state.UserFind(data.Login)
			if user == nil || !user.CheckPassword(data.Password) {
				return nil, jwt.ErrFailedAuthentication
			}
			log.Println("[DEBUG]: Authenticator2")

			return &User{
				Login: user.Login,
				Name:  user.Name,
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			log.Println("[DEBUG]: Authorizator", data.(*User))
			// TODO: check the func
			if v, ok := data.(*User); ok && v.Login == "admin" {
				return true
			}

			return false
		},
	})

	if err != nil {
		log.Panic("JWT Error:", err)
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
		user := v1.Group("/user")
		user.Use(api_data.JWT.MiddlewareFunc())
		{
			user.GET("/", UserGetList)
			user.GET("/:login", UserGet)
			user.POST("/:login", UserPost)
			user.DELETE("/:login", UserDelete)
		}
		access := v1.Group("/access")
		access.Use(api_data.JWT.MiddlewareFunc())
		{
			access.GET("/", AccessGetList)
			access.GET("/:id", AccessGet)
			access.POST("/*id", AccessPost)
			access.DELETE("/:id", AccessDelete)
		}
		source := v1.Group("/source")
		source.Use(api_data.JWT.MiddlewareFunc())
		{
			source.GET("/", SourceGetList)
			source.GET("/:id", SourceGet)
			source.POST("/:id", SourcePost)
			source.DELETE("/:id", SourceDelete)
		}
		navigate := v1.Group("/navigate")
		navigate.Use(api_data.JWT.MiddlewareFunc())
		{
			navigate.GET("/*path", NavigateGetList)
		}
	}
}

type api_s struct {
	JWT *jwt.GinJWTMiddleware // jwt object
}

var api_data = &api_s{}
