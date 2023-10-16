package handlers

import (
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	mariadb "suglider-auth/internal/database"
	"suglider-auth/pkg/encrypt"
	"database/sql"
	"suglider-auth/pkg/session"
)

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
				c.Redirect(http.StatusMovedPermanently, "/hello")
			} else if !pwdVerify {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				return
			}

		} else if err == sql.ErrNoRows {
			log.Println("User Login failed:", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		} else if err != nil {
			log.Println("Login failed:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
	} else {
		c.Redirect(http.StatusMovedPermanently, "/hello")
	}
}

func UserLoginForm(c *gin.Context) {
	err := c.Request.ParseMultipartForm(2048)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "FormData Invalid"})
		return
	}
	username := c.PostForm("username")
	password := c.PostForm("password")
	log.Println("Username: ", username)
	log.Println("Password: ", password)
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}
	} else if err == sql.ErrNoRows {
		log.Println("User Login failed:", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	} else if err != nil {
		log.Println("Login failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
}

func UserLogOutDemo(c *gin.Context) {
	sid := session.ReadSession(c)
	ok := session.CheckSession(c)
	if !ok {
		log.Printf("session ID %s doesn't exsit in redis\n", sid)
		return
	}
	session.DeleteSession(sid)
	c.Redirect(http.StatusMovedPermanently, "/login")
}
