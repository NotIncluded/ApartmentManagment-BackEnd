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

	// Validate required fields are not empty
	if loginRequest.Email == "" || loginRequest.Password == "" {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Incomplete request body: email and password are required",
			nil,
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
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
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

	// Validate required fields are not empty
	if req.Name == "" || req.Phone == "" || req.Email == "" || req.Password == "" {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Incomplete request body: name, phone, email, and password are required",
			nil,
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	// Reject if trying to set ADMIN role during registration
	if req.Role == "ADMIN" {
		appErr := response.NewAppResponse(
			http.StatusBadRequest,
			"Cannot set ADMIN role during tenant registration",
			nil,
		)
		c.JSON(appErr.Status, appErr.Response())
		return
	}

	createdUser, err := ac.authService.Register(req.Name, req.Phone, req.Email, req.Password, "TENANT")
	if err != nil {
		appErr := response.NewAppResponse(http.StatusBadRequest, "Registration failed", err.Error())
		c.JSON(appErr.Status, appErr.Response())
		return
	}
	res := response.NewAppResponse(http.StatusOK, "User registered successfully", nil)
	res.Data = createdUser
	c.JSON(res.Status, res.Response())
}
