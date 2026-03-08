package controller

import (
	"fmt"
	"net/http"

	"github.com/PunMung-66/ApartmentSys/internal/response"
	"github.com/PunMung-66/ApartmentSys/model"
	"github.com/PunMung-66/ApartmentSys/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (u *UserController) CreateUser(c *gin.Context) {

	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		res := response.NewAppResponse(
			http.StatusBadRequest,
			"Invalid request body",
			err.Error(),
		)
		c.JSON(res.Status, res.Response())
		return
	}

	userResponse, err := u.userService.CreateUser(&user)

	if err != nil {
		res := response.NewAppResponse(
			http.StatusBadRequest,
			err.Error(),
			nil,
		)
		c.JSON(res.Status, res.Response())
		return
	}

	res := response.NewAppResponse(
		http.StatusAccepted,
		"Add user successfully",
		userResponse,
	)

	c.JSON(res.Status, res.Response())
}

func (u *UserController) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	var id string
	if _, err := fmt.Sscanf(userID, "%s", &id); err != nil {
		res := response.NewAppResponse(
			http.StatusBadRequest,
			"Invalid user ID",
			err.Error(),
		)
		c.JSON(res.Status, res.Response())
		return
	}

	err := u.userService.DeleteUser(id)

	if err != nil {
		res := response.NewAppResponse(
			http.StatusBadRequest,
			"Failed to delete user",
			err.Error(),
		)
		c.JSON(res.Status, res.Response())
		return
	}

	res := response.NewAppResponse(
		http.StatusOK,
		"User deleted successfully",
		nil,
	)

	c.JSON(res.Status, res.Response())
}
