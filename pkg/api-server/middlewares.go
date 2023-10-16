package api_server

import (
	"log/slog"
	"net/http"
	"github.com/gin-gonic/gin"
	"suglider-auth/pkg/rbac"
)

func userPrivilege(csbn *rbac.CasbinEnforcerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		sub, exist := c.Get("Username")
		if !exist {
			sub = "anonymous"
		}
		obj := c.Request.URL.Path  // c.Request.URL.RequestURI()
		act := c.Request.Method

		switch c.Request.URL.Path {
		case "/login", "/sing-up", "/api/v1/user/login",
			 "/api/v1/user/logout", "/api/v1/user/sing-up":
				sub = "anonymous"
		default:
		}

		if pass, err := csbn.Enforcer.Enforce(sub.(string), obj, act); !pass {
			if err != nil {
				slog.ErrorContext(c, err.Error())
			}
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
			// c.AbortWithStatusJSON(http.StatusForbidden, gin.H {
			// 	"status": "forbidden",
			// 	"message": "You don't have the required permission.",
			// })
		}
		c.Next()
	}
}
