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
		ss := session.ReadSession(c)
		info := session.SessionData{}
		if err := json.Unmarshal([]byte(ss), &info); err != nil {
			slog.Error(err.Error())
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
		}
		c.Set("Username", info.Username)
		c.Next()
	}
}
