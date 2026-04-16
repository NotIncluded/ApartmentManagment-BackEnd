package controller

import (
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
		http.StatusCreated,
		"Add user successfully",
		userResponse,
	)

	c.JSON(res.Status, res.Response())
}

func (u *UserController) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	role, _ := c.Get("role")
	currentUserID, _ := c.Get("user_id")

	if role == "TENANT" && currentUserID != userID {
		res := response.NewAppResponse(
			http.StatusForbidden,
			"You can only view your own profile",
			nil,
		)
		c.JSON(res.Status, res.Response())
		return
	}

	userResponse, err := u.userService.GetUserByID(userID)
	if err != nil {
		res := response.NewAppResponse(
			http.StatusNotFound,
			"User not found",
			err.Error(),
		)
		c.JSON(res.Status, res.Response())
		return
	}

	res := response.NewAppResponse(
		http.StatusOK,
		"User retrieved successfully",
		userResponse,
	)
	c.JSON(res.Status, res.Response())
}

func (u *UserController) GetUsersByRole(c *gin.Context) {
	role := c.Query("role")

	if role == "" {
		res := response.NewAppResponse(
			http.StatusBadRequest,
			"Role query parameter is required",
			nil,
		)
		c.JSON(res.Status, res.Response())
		return
	}

	users, err := u.userService.GetUsersByRole(role)
	if err != nil {
		res := response.NewAppResponse(
			http.StatusInternalServerError,
			"Failed to retrieve users",
			err.Error(),
		)
		c.JSON(res.Status, res.Response())
		return
	}

	res := response.NewAppResponse(
		http.StatusOK,
		"Users retrieved successfully",
		users,
	)
	c.JSON(res.Status, res.Response())
}

func (u *UserController) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	role, _ := c.Get("role")
	currentUserID, _ := c.Get("user_id")

	if role == "TENANT" && currentUserID != userID {
		res := response.NewAppResponse(
			http.StatusForbidden,
			"You can only update your own profile",
			nil,
		)
		c.JSON(res.Status, res.Response())
		return
	}

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

	user.ID = userID
	if role == "TENANT" {
		user.Role = "TENANT"
	}

	userResponse, err := u.userService.UpdateUser(&user)
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
		http.StatusOK,
		"User updated successfully",
		userResponse,
	)
	c.JSON(res.Status, res.Response())
}

func (u *UserController) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	err := u.userService.DeleteUser(userID)
	if err != nil {
		res := response.NewAppResponse(
			http.StatusNotFound,
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
