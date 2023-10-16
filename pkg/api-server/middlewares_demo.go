package api_server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"github.com/gin-gonic/gin"
	"suglider-auth/pkg/session"
)

func checkSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		ss, err := session.ReadSessionData(c)
		if err != nil {
			c.Set("Username", "anonymous")
		} else {
			info := session.SessionData{}
			if err := json.Unmarshal([]byte(ss), &info); err != nil {
					slog.Error(err.Error())
					c.Redirect(http.StatusTemporaryRedirect, "/login")
					c.Abort()
			}
			c.Set("Username", info.Username)
			switch c.Request.URL.Path {
			case "/login", "/api/v1/login":
				c.Redirect(http.StatusTemporaryRedirect, "/hello")
			default:
			}
		}
		c.Next()
	}
}
