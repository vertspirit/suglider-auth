package session

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"suglider-auth/internal/redis"
)

type SessionData struct {
	Username    string    `json:"username"`
}

func ReadSessionData(c *gin.Context) (string, error) {
	session := sessions.Default(c)
	id := session.Get("sid")
	if id == nil {
		return "", fmt.Errorf("Session ID not found.\n")
	}
	ss, err := redis.GetCtx(c, fmt.Sprintf("sid:%v", id))
	if err != nil {
		return "", err
	}
	return ss, nil
}
