package user

import (
	"github.com/gin-gonic/gin"
	"suglider-auth/pkg/api-server/api_v1/handlers"
)

func UserHandler(router *gin.RouterGroup) {

	// demo web
	router.POST("/sign-up", handlers.UserSignUpDemo)
	router.POST("/delete", handlers.UserDelete)
	router.POST("/login", handlers.UserLoginDemo)
	router.POST("/logout", handlers.UserLogOutDemo)

	// Test
	router.GET("/test", handlers.Test)
	router.GET("/test-v2", handlers.Testv2)
}
