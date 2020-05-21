package api

import (
	"log"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/state-of-the-art/NyanSync/lib/state"
)

type jwt_authenticate_body struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func initAuthV1() {
	identity_key := "id"
	var login_data jwt_authenticate_body

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
			if c.ShouldBind(&login_data) != nil {
				return "", jwt.ErrMissingLoginValues
			}

			user := state.UserFind(login_data.Login)
			if user == nil || !user.CheckPassword(login_data.Password) {
				return nil, jwt.ErrFailedAuthentication
			}

			return &User{
				Login: user.Login,
				Name:  user.Name,
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			log.Println("[DEBUG]: Authorizator", data.(*User))
			// TODO: check the policies of the user
			if v, ok := data.(*User); ok && v.Login == "admin" {
				return true
			}

			return false
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			// Getting the user, stored during login Authenticate
			user := state.UserFind(login_data.Login)
			c.JSON(http.StatusOK, gin.H{
				"code":   http.StatusOK,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
				"user":   toAPIUser(user),
			})
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
			// TODO: Make sure the token can't be used anymore
			auth.POST("/logout", api_data.JWT.LogoutHandler)
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
