package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"github.com/gin-gonic/gin"
	mariadb "suglider-auth/internal/database"
	"suglider-auth/pkg/encrypt"
	"database/sql"
	"suglider-auth/pkg/session"
)

func UserSignUpDemo(c *gin.Context) {
	var request userSignUp
	var err error

	// Check the parameter trasnfer from POST
	err = c.ShouldBindJSON(&request)
	if err != nil {
		c.HTML(http.StatusBadRequest, "signup.html", gin.H {"error": err.Error()})
		return
	}

	// Encode user password
	passwordEncode, _ := encrypt.SaltedPasswordHash(request.Password)

	err = mariadb.UserSignUp(request.Username, passwordEncode, request.Mail, request.Address)
	if err != nil {
		slog.ErrorContext(c, err.Error())
		c.HTML(http.StatusInternalServerError, "signup.html", gin.H {"error": "Internal Server Error"})
		return
	} else {
		// c.JSON(http.StatusOK, gin.H {"message": "User created successfully"})
		c.Redirect(http.StatusMovedPermanently, "/login")
	}
}

func UserLoginDemo(c *gin.Context) {

	var request userLogin
	// var userDBInfo userDBInfo

	// Check the parameter trasnfer from POST
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ok := session.CheckSession(c)
	if !ok {
		// Check whether username exist or not
		userInfo, err := mariadb.UserLogin(request.Username)

		// No err means user exist
		if err == nil {

			pwdVerify := encrypt.VerifySaltedPasswordHash(userInfo.Password, request.Password)

			// Check password true or false
			if pwdVerify {
				// c.JSON(http.StatusOK, gin.H{"message": "User Logined successfully"})
				session.AddSession(c, request.Username)
			} else if !pwdVerify {
				c.HTML(http.StatusUnauthorized, "login.html", gin.H {"error": "Invalid password"})
				return
			} else {
				c.HTML(http.StatusInternalServerError, "login.html", gin.H {"error": "Internal Server Error"})
				return
			}

		} else if err == sql.ErrNoRows {
			slog.ErrorContext(c, err.Error())
			c.HTML(http.StatusNotFound, "login.html", gin.H {"error": "User not found"})
			return
		} else if err != nil {
			slog.ErrorContext(c, err.Error())
			c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Internal Server Error"})
			return
		}
	}
	c.Redirect(http.StatusMovedPermanently, "/hello")
}

func UserLoginForm(c *gin.Context) {
	err := c.Request.ParseMultipartForm(2048)
	if err != nil {
		slog.ErrorContext(c, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "FormData Invalid"})
		return
	}
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username Invalid"})
		return
	}

	userInfo, err := mariadb.UserLogin(username)
	if err == nil {
		pwdVerify := encrypt.VerifySaltedPasswordHash(userInfo.Password, password)

		// Check password true or false
		if pwdVerify {
			session.AddSession(c, username)
			c.Redirect(http.StatusMovedPermanently, "/hello")
		} else if !pwdVerify {
			c.HTML(http.StatusUnauthorized, "login.html", gin.H {"error": "Invalid password"})
			return
		} else {
			c.HTML(http.StatusInternalServerError, "login.html", gin.H {"error": "Internal Server Error"})
			return
		}
	} else if err == sql.ErrNoRows {
		slog.ErrorContext(c, err.Error())
		c.HTML(http.StatusNotFound, "login.html", gin.H {"error": "User not found"})
		return
	} else if err != nil {
		slog.ErrorContext(c, err.Error())
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"error": "Internal Server Error"})
		return
	}
}

func UserLogOutDemo(c *gin.Context) {
	sid := session.ReadSession(c)
	ok := session.CheckSession(c)
	if !ok {
		slog.ErrorContext(c, fmt.Sprintf("session ID %s doesn't exsit in redis\n", sid))
	}
	session.DeleteSession(sid)
	c.Redirect(http.StatusMovedPermanently, "/login")
}
