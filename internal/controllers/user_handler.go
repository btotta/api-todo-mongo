package controllers

import (
	"todo-app-mongo/internal/database"
	"todo-app-mongo/internal/middleware"
	"todo-app-mongo/internal/models/dtos"
	"todo-app-mongo/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserHandlerInterface interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetUser(c *gin.Context)
	Login(c *gin.Context)
	Refresh(c *gin.Context)
	Logout(c *gin.Context)
}

type userHandler struct {
	userDAO        database.UserDAOInterface
	authMiddleware middleware.AuthMiddlewareInterface
}

func NewUserHandler(userDAO database.UserDAOInterface, auth middleware.AuthMiddlewareInterface) UserHandlerInterface {
	return &userHandler{userDAO: userDAO,
		authMiddleware: auth}

}

// @Summary Create a new user
// @Description Create a new user
// @Tags user
// @Accept json
// @Produce json
// @Param user body dtos.UserRequestDTO true "User object"
// @Success 201 {object} dtos.UserResponseDTO "User created"
// @Failure 400 {object} utils.ErrorHandler
// @Router /user [post]
func (u *userHandler) Create(c *gin.Context) {

	var user dtos.UserRequestDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.DefaultErrorResponse(c, 400, "Invalid request body")
		return
	}

	userModel, err := user.ToUserModel()
	if err != nil {
		utils.DefaultErrorResponse(c, 400, err.Error())
		return
	}

	err = u.userDAO.Create(c, userModel)
	if err != nil {
		utils.DefaultErrorResponse(c, 500, "Internal server error")
		return
	}

	c.JSON(201, dtos.UserResponseDTO{
		Name:  userModel.Name,
		Email: userModel.Email,
		ID:    userModel.ID.Hex(),
	})
}

// @Summary Login
// @Description Login
// @Tags user
// @Accept json
// @Produce json
// @Param user body dtos.UserLoginDTO true "User object"
// @Success 200 {object} dtos.UserLoginResponseDTO "User logged in"
// @Failure 400 {object} utils.ErrorHandler
// @Router /login [post]
func (u *userHandler) Login(c *gin.Context) {

	var dto dtos.UserLoginDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.DefaultErrorResponse(c, 400, "Invalid request body")
		return
	}

	user, err := u.userDAO.GetByEmail(c, dto.Email)
	if err != nil {
		utils.DefaultErrorResponse(c, 400, "Invalid email or password")
		return
	}

	if !user.ComparePassword(dto.Password) {
		utils.DefaultErrorResponse(c, 400, "Invalid email or password")
		return
	}

	tokenResponse, err := u.authMiddleware.GenerateTokenAndRefresh(user)
	if err != nil {
		utils.DefaultErrorResponse(c, 500, "Error generating tokens")
		return
	}

	c.JSON(200, tokenResponse)
}

// @Summary Get user by ID
// @Description Get user by ID
// @Security Bearer
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} dtos.UserResponseDTO "User found"
// @Failure 404 {object} utils.ErrorHandler
// @Router /user [get]
func (u *userHandler) GetUser(c *gin.Context) {
	// pega email do user do contexto
	email := c.GetString("email")

	user, err := u.userDAO.GetByEmail(c, email)
	if err != nil {
		utils.DefaultErrorResponse(c, 404, "User not found")
		return
	}

	c.JSON(200, dtos.UserResponseDTO{
		ID:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
	})
}

// @Summary Update user
// @Description Update user
// @Security Bearer
// @Tags user
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body dtos.UserRequestDTO true "User object"
// @Success 200 {object} dtos.UserResponseDTO "User updated"
// @Failure 400 {object} utils.ErrorHandler
// @Router /user [put]
func (u *userHandler) Update(c *gin.Context) {
	email := c.GetString("email")

	var user dtos.UserRequestDTO
	if err := c.ShouldBindJSON(&user); err != nil {
		utils.DefaultErrorResponse(c, 400, "Invalid request body")
		return
	}

	dbUser, err := u.userDAO.GetByEmail(c, email)
	if err != nil {
		utils.DefaultErrorResponse(c, 404, "User not found")
		return
	}

	// por enquanto só atualiza o nome
	dbUser.Name = user.Name

	err = u.userDAO.Update(c, dbUser)
	if err != nil {
		utils.DefaultErrorResponse(c, 500, "Internal server error")
		return
	}

	c.JSON(200, dtos.UserResponseDTO{
		Name:  dbUser.Name,
		Email: dbUser.Email,
		ID:    dbUser.ID.Hex(),
	})
}

// @Summary Logout
// @Description Logout
// @Security Bearer
// @Tags user
// @Accept json
// @Produce json
// @Success 200
// @Router /logout [post]
func (u *userHandler) Logout(c *gin.Context) {
	loggedOut := u.authMiddleware.LogOff(c)
	if !loggedOut {
		return
	}
	c.JSON(200, gin.H{
		"message": "Logged out successfully",
		"succes":  true,
	})
}

// @Summary Delete user
// @Description Delete user
// @Security Bearer
// @Tags user
// @Accept json
// @Produce json
// @Success 200
// @Router /user [delete]
func (u *userHandler) Delete(c *gin.Context) {
	email := c.GetString("email")

	err := u.userDAO.Delete(c, email)
	if err != nil {
		utils.DefaultErrorResponse(c, 500, "Internal server error")
		return
	}

	c.JSON(200, gin.H{
		"message": "User deleted",
		"success": true,
	})
}

// @Summary Refresh token
// @Description Refresh token
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} dtos.UserLoginResponseDTO "Token refreshed"
// @Failure 400 {object} utils.ErrorHandler
// @Router /refresh [post]
func (u *userHandler) Refresh(c *gin.Context) {
	tokenResponse := u.authMiddleware.RefreshToken(c)
	if tokenResponse == nil {
		return
	}

	c.JSON(200, tokenResponse)
}
