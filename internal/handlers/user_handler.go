package handlers

import (
	"todo-app-mongo/internal/database"
	"todo-app-mongo/internal/dtos"
	"todo-app-mongo/internal/pkg/security"
	"todo-app-mongo/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userDAO database.UserDAOInterface
}

func NewUserHandler(userDAO database.UserDAOInterface) *UserHandler {
	return &UserHandler{userDAO: userDAO}

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
func (u *UserHandler) Create(c *gin.Context) {

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

	userModel, err = u.userDAO.Create(c, userModel)
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
func (u *UserHandler) Login(c *gin.Context) {

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

	token, err := security.GenerateToken(user.Email)
	if err != nil {
		utils.DefaultErrorResponse(c, 500, "Internal server error")
		return
	}

	refreshToken, err := security.GenerateRefreshToken(user.Email)
	if err != nil {
		utils.DefaultErrorResponse(c, 500, "Internal server error")
		return
	}

	c.JSON(200, dtos.UserLoginResponseDTO{
		Token:        token,
		RefreshToken: refreshToken,
	})

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
func (u *UserHandler) GetUser(c *gin.Context) {
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
func (u *UserHandler) Update(c *gin.Context) {
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

	dbUser, err = u.userDAO.Update(c, dbUser)
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
func (u *UserHandler) Logout(c *gin.Context) {

	token := c.GetHeader("Authorization")
	refreshtoken := c.GetHeader("Refresh")

	if token == "" || refreshtoken == "" {
		utils.DefaultErrorResponse(c, 400, "Invalid request")
	}

	security.LogOff(token)
	security.LogOff(refreshtoken)

	if !security.IsLoggedOff(token) && !security.IsLoggedOff(refreshtoken) {
		utils.DefaultErrorResponse(c, 500, "Internal server error")
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
func (u *UserHandler) Delete(c *gin.Context) {
	email := c.GetString("email")

	_, err := u.userDAO.Delete(c, email)
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
func (u *UserHandler) Refresh(c *gin.Context) {

	token := c.GetHeader("Authorization")
	refreshtoken := c.GetHeader("Refresh")

	if token == "" || refreshtoken == "" {
		utils.DefaultErrorResponse(c, 400, "Invalid request")
	}

	email, err := security.ValidateRefreshToken(refreshtoken)
	if err != nil {
		utils.DefaultErrorResponse(c, 401, "Invalid refresh token")
		return
	}

	if !security.IsLoggedOff(token) {
		utils.DefaultErrorResponse(c, 401, "Invalid token")
		return
	}

	newToken, err := security.GenerateToken(email)
	if err != nil {
		utils.DefaultErrorResponse(c, 500, "Internal server error")
		return
	}

	c.JSON(200, dtos.UserLoginResponseDTO{
		Token:        newToken,
		RefreshToken: refreshtoken,
	})

}
