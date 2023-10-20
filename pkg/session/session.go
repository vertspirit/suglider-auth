package session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"suglider-auth/pkg/encrypt"
	"suglider-auth/internal/redis"
	"encoding/json"
	"log/slog"
	"suglider-auth/pkg/time_convert"
	"fmt"
  )

type sessionData struct {
	Username	string	`json:"username"`
}

func AddSession(c *gin.Context, user string) (string, error) {

	sessionValue := sessionData{
		Username: user,
	}

	jsonSessionValue, err := json.Marshal(sessionValue)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to create session value JSON data: %v", err)
		slog.Error(errorMessage)
		return "", err
	}

	// Genertate session ID with no Dash
	sessionID := encrypt.GenertateUUID(true)

	session := sessions.Default(c)

	session.Set("sid", sessionID)
	session.Save()

	redisKey := "sid:" + sessionID
	redisValue := string(jsonSessionValue)

	// time_convert.RedisTTL is a global variable from time_convert.go
	redis.Set(redisKey, redisValue, time_convert.RedisTTL)
	
	return sessionID, nil
}

func CheckSession(c *gin.Context) bool {
	session := sessions.Default(c)
	sid := session.Get("sid")

	// Process key format
	sessionKey := fmt.Sprintf("sid:%s", sid)
	
	// Exists() function will return bool
	return redis.Exists(sessionKey)
}

func DeleteSession(sid string) {
	// Process key format
	sessionKey := fmt.Sprintf("sid:%s", sid)
	redis.Delete(sessionKey)
}

func ReadSession(c *gin.Context) string {
	session := sessions.Default(c)
	sid := session.Get("sid")

	// Convert interface{} type to String type
	StrSid := fmt.Sprintf("%v", sid)
	return StrSid
}