package api_server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HelloPage(c *gin.Context){
	c.HTML(http.StatusOK, "index.html", nil)
}

func LoginPage(c *gin.Context){
	c.HTML(http.StatusOK, "login.html", nil)
}

func Login2Page(c *gin.Context){
	c.HTML(http.StatusOK, "login_form", nil)
}

func SignUpPage(c *gin.Context){
	c.HTML(http.StatusOK, "signup.html", nil)
}
