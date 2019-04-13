package http

import (
	"ginexamples"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *AppServer) RegisterUserHandler(c *gin.Context) {
	type request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var (
		userModel ginexamples.User
		req       request
	)

	err := c.BindJSON(&req)
	if err != nil || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userModel.Email = req.Email
	userModel.Name = req.Name

	user, err := a.UserService.CreateUser(&userModel, req.Password)
	if err != nil {
		a.Logger.Printf("error creating user: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setCookie(c, user.SessionID)
	c.JSON(http.StatusOK, gin.H{
		"ID":    user.ID,
		"Name":  user.Name,
		"Email": user.Email,
	})
}

func (a *AppServer) LoginUserHandler(c *gin.Context) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	var req request
	err := c.BindJSON(&req)
	if err != nil || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := a.UserService.Login(req.Email, req.Password)
	if err != nil {
		a.Logger.Printf("error logging in: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	setCookie(c, user.SessionID)
	c.JSON(http.StatusOK, gin.H{
		"Name":  user.Name,
		"Email": user.Email,
	})
}

func setCookie(c *gin.Context, value string) {
	c.SetCookie("sessionID", value, 86400, "/", "localhost", false, true)
}
