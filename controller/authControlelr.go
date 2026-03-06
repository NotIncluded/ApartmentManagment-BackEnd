package controller

import (
	"net/http"

	"github.com/PunMung-66/ApartmentSys/internal/response"
	"github.com/PunMung-66/ApartmentSys/service"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
	signature   []byte
}

type AuthControllerInterface interface {
	LoginHandler(c *gin.Context)
	RegisterHandler(c *gin.Context)
}

func NewAuthController(authService *service.AuthService, signature []byte) *AuthController {
	return &AuthController{
		authService: authService,
		signature:   signature,
	}
}

func (ac *AuthController) LoginHandler(c *gin.Context) {
	var loginRequest service.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	tokenString, err := ac.authService.Login(loginRequest, ac.signature)
	if err != nil {
		appErr := response.NewAppResponse(http.StatusUnauthorized, "Invalid username or password", err.Error())
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	res := response.NewAppResponse(http.StatusOK, "Access token generated successfully", nil)
	res.Data = map[string]string{
		"access_token": tokenString,
	}
	c.JSON(res.Status, res.Response())
}

func (ac *AuthController) RegisterHandler(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	createdUser, err := ac.authService.Register(req.Username, req.Password)
	if err != nil {
		appErr := response.NewAppResponse(http.StatusBadRequest, "Registration failed", err.Error())
		c.JSON(appErr.Status, appErr.Response())
		return
	}
	res := response.NewAppResponse(http.StatusOK, "User registered successfully", nil)
	res.Data = createdUser
	c.JSON(res.Status, res.Response())
}