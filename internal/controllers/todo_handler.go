package controllers

import (
	"errors"
	"strconv"
	"todo-app-mongo/internal/database"
	"todo-app-mongo/internal/models"
	"todo-app-mongo/internal/models/dtos"
	"todo-app-mongo/internal/utils"

	"github.com/gin-gonic/gin"
)

type TodoHandlerInterface interface {
	Create(c *gin.Context)
	Get(c *gin.Context)
	GetAll(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type todoHandler struct {
	todoDAO database.TodoDAOInterface
	userDAO database.UserDAOInterface
}

func NewTodoHandler(todoDAO database.TodoDAOInterface, userDAO database.UserDAOInterface) TodoHandlerInterface {
	return &todoHandler{todoDAO: todoDAO, userDAO: userDAO}
}

// @Summary Create a new todo
// @Description Create a new todo
// @Tags todo
// @Accept json
// @Produce json
// @Param todo body dtos.TodoDTO true "Todo object"
// @Success 201 {object} models.Todo
// @Failure 400 {object} utils.ErrorHandler
// @Router /todo [post]
func (t *todoHandler) Create(c *gin.Context) {

	user, err := t.getUserFromContext(c)
	if err != nil {
		utils.DefaultErrorResponse(c, 500, "Error getting user")
		return
	}

	var todoDTO dtos.TodoDTO
	if err := c.ShouldBindJSON(&todoDTO); err != nil {
		utils.DefaultErrorResponse(c, 400, "Invalid request body")
		return
	}

	todo := todoDTO.ToModel()
	todo.UserID = user.ID

	if err := t.todoDAO.Create(c, todo); err != nil {
		utils.DefaultErrorResponse(c, 500, "Error creating todo")
		return
	}

	c.JSON(201, todo)

}

// @Summary Get a todo by ID
// @Description Get a todo by ID
// @Tags todo
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Success 200 {object} models.Todo
// @Failure 500 {object} utils.ErrorHandler
// @Router /todo/{id} [get]
func (t *todoHandler) Get(c *gin.Context) {

	user, err := t.getUserFromContext(c)
	if err != nil {
		utils.DefaultErrorResponse(c, 500, "Error getting user")
		return
	}

	id := c.Param("id")
	todo, err := t.todoDAO.Get(c, id, user.ID.String())
	if err != nil {
		utils.DefaultErrorResponse(c, 500, "Error getting todo")
		return
	}

	c.JSON(200, todo)
}

// @Summary Get all todos
// @Description Get all todos
// @Tags todo
// @Accept json
// @Produce json
// @Success 200 {array} models.Todo
// @Failure 500 {object} utils.ErrorHandler
// @Router /todos [get]
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
func (t *todoHandler) GetAll(c *gin.Context) {

	var limit int64
	var offset int64
	var search string

	l := c.Query("limit")
	if l == "" {
		limit = 10
	} else {
		limit, _ = strconv.ParseInt(l, 10, 64)
	}

	o := c.Query("offset")
	if o == "" {
		offset = 0
	} else {
		offset, _ = strconv.ParseInt(o, 10, 64)
	}

	s := c.Query("search")
	if s != "" {
		search = s
	}

	todos, count, err := t.todoDAO.GetAll(c, limit, offset, search)
	if err != nil {
		utils.DefaultErrorResponse(c, 500, "Error getting todos")
		return
	}

	c.JSON(200, dtos.NewPageDTO(int64(len(todos)), offset, count, todos))

}

// @Summary Update a todo by ID
// @Description Update a todo by ID
// @Tags todo
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Param todo body dtos.TodoDTO true "Todo object"
// @Success 204 {object} models.Todo
// @Failure 500 {object} utils.ErrorHandler
// @Router /todo/{id} [put]
func (t *todoHandler) Update(c *gin.Context) {

	id := c.Param("id")
	var todoDTO dtos.TodoDTO
	if err := c.ShouldBindJSON(&todoDTO); err != nil {
		utils.DefaultErrorResponse(c, 400, "Invalid request body")
		return
	}

	todo := todoDTO.ToModel()
	if err := t.todoDAO.Update(c, id, todo); err != nil {
		utils.DefaultErrorResponse(c, 500, "Error updating todo")
		return
	}

	c.JSON(204, todo)
}

// @Summary Delete a todo by ID
// @Description Delete a todo by ID
// @Tags todo
// @Accept json
// @Param id path string true "Todo ID"
// @Success 204
// @Failure 500 {object} utils.ErrorHandler
// @Router /todo/{id} [delete]
func (t *todoHandler) Delete(c *gin.Context) {

	id := c.Param("id")
	if err := t.todoDAO.Delete(c, id); err != nil {
		utils.DefaultErrorResponse(c, 500, "Error deleting todo")
		return
	}

	c.JSON(204, nil)
}

func (t *todoHandler) getUserFromContext(c *gin.Context) (*models.User, error) {

	user, err := t.userDAO.GetByEmail(c, c.GetString("email"))
	if err != nil {
		return nil, errors.New("error getting user")
	}

	return user, nil
}
