package api

import (
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/state-of-the-art/NyanSync/lib/rbac"
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
			if state.UserExists(claims[identity_key].(string)) {
				user := state.UserGet(claims[identity_key].(string))
				return user
			}
			return nil
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			log.Println("[DEBUG]: Authenticator")
			if c.ShouldBind(&login_data) != nil {
				return "", jwt.ErrMissingLoginValues
			}

			user := state.UserGet(login_data.Login)
			if !user.PasswordCheck(login_data.Password) {
				return nil, jwt.ErrFailedAuthentication
			}

			return &User{
				Login: user.Login,
				Name:  user.Name,
			}, nil
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			log.Println("[DEBUG]: Authorizator")
			_, ok := data.(*state.User)
			if !ok {
				return false
			}
			return true
		},
		LoginResponse: func(c *gin.Context, code int, token string, expire time.Time) {
			// Getting the user, stored during login Authenticate
			user := state.UserGet(login_data.Login)
			c.JSON(http.StatusOK, gin.H{
				"code":   http.StatusOK,
				"token":  token,
				"expire": expire.Format(time.RFC3339),
				"user":   toAPIUser(&user),
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
	v1.Use(
		// Skip auth to allow unknown visitor to login
		func(c *gin.Context) {
			// Allow api/auth to be used without authorization
			if strings.HasPrefix(c.FullPath(), v1.BasePath()+"/auth") {
				c.Handler()(c)
				c.Abort()
				return
			}
			c.Next()
		},
		// Processing JWT
		api_data.JWT.MiddlewareFunc(),
		// Processing RBAC
		ProcessRBAC,
		// TODO: check access to navigator path
	)
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", api_data.JWT.LoginHandler)
			// TODO: Make sure the token can't be used anymore
			auth.POST("/logout", api_data.JWT.LogoutHandler)
			auth.GET("/refresh_token", api_data.JWT.RefreshHandler)
		}
		user := v1.Group("/user")
		{
			user.GET("/", UserGetList)
			user.GET("/:login", UserGet)
			user.POST("/:login", UserPost)
			user.DELETE("/:login", UserDelete)
		}
		access := v1.Group("/access")
		{
			access.GET("/", AccessGetList)
			access.GET("/:id", AccessGet)
			access.POST("/*id", AccessPost)
			access.DELETE("/:id", AccessDelete)
		}
		source := v1.Group("/source")
		{
			source.GET("/", SourceGetList)
			source.GET("/:id", SourceGet)
			source.POST("/:id", SourcePost)
			source.DELETE("/:id", SourceDelete)
		}
		navigate := v1.Group("/navigate")
		{
			navigate.GET("/*path", NavigateGetList)
		}
	}

	// Prepare permissions list
	r := state.GetRBAC()
	for _, route := range router.Routes() {
		if !strings.HasPrefix(route.Path, v1.BasePath()) || strings.HasPrefix(route.Path, v1.BasePath()+"/auth") {
			continue
		}

		if rbac.PermSet(r, getPermId(route.Path), route.Method) {
			// TODO: simplify rbac save
			state.SaveRBAC()
		}

	}
}

func ContextAccount(c *gin.Context) *state.User {
	data, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Account is not set"})
		log.Panic("Critical issue during getting account from context")
	}
	acc, ok := data.(*state.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Wrong account type in context"})
		log.Panic("Critical issue during getting account from context")
	}

	return acc
}

func ProcessRBAC(c *gin.Context) {
	acc := ContextAccount(c)
	if acc.Role == "" {
		c.JSON(http.StatusForbidden, gin.H{"message": "Account role is not set"})
		c.Abort()
		return
	}

	r := state.GetRBAC()
	perm := r.GetPermission(getPermId(c.FullPath()))
	if perm == nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unable to find RBAC permission"})
		c.Abort()
		return
	}

	// Check permissions
	switch c.Request.Method {
	case "GET":
		if !r.IsGranted(acc.Role, perm, rbac.Read) {
			c.JSON(http.StatusForbidden, gin.H{"message": "No read access"})
			c.Abort()
		}
	case "POST":
		if !r.IsGranted(acc.Role, perm, rbac.Update) ||
			!r.IsGranted(acc.Role, perm, rbac.Create) ||
			!r.IsGranted(acc.Role, perm, rbac.UpdateSelf) {
			c.JSON(http.StatusForbidden, gin.H{"message": "No update/create access"})
			c.Abort()
		}
	case "DELETE":
		if !r.IsGranted(acc.Role, perm, rbac.Delete) {
			c.JSON(http.StatusForbidden, gin.H{"message": "No delete access"})
			c.Abort()
		}
	default:
		c.JSON(http.StatusForbidden, gin.H{"message": "Unknown method"})
		c.Abort()
	}

	c.Next()
}

func getPermId(full_path string) string {
	split := strings.Split(full_path, "/")
	split = split[0 : len(split)-1]
	return strings.Join(split, "/")
}

type api_s struct {
	JWT *jwt.GinJWTMiddleware // jwt object
}

var api_data = &api_s{}
