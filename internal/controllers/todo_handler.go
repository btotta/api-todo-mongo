package controllers

import (
	"strconv"
	"todo-app-mongo/internal/database"
	"todo-app-mongo/internal/models/dtos"

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
}

func NewTodoHandler(todoDAO database.TodoDAOInterface) TodoHandlerInterface {
	return &todoHandler{todoDAO: todoDAO}
}

// @Summary Create a new todo
// @Description Create a new todo
// @Tags todo
// @Accept json
// @Produce json
// @Param todo body dtos.TodoDTO true "Todo object"
// @Success 201 {object} models.Todo
// @Router /todo [post]
func (t *todoHandler) Create(c *gin.Context) {

	var todoDTO dtos.TodoDTO
	if err := c.ShouldBindJSON(&todoDTO); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	todo := todoDTO.ToModel()
	if err := t.todoDAO.Create(c, todo); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
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
// @Router /todo/{id} [get]
func (t *todoHandler) Get(c *gin.Context) {

	id := c.Param("id")
	todo, err := t.todoDAO.Get(c, id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
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
// @Router /todos [get]
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
func (t *todoHandler) GetAll(c *gin.Context) {

	var limit int64
	var offset int64

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

	todos, err := t.todoDAO.GetAll(c, limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, dtos.NewPageDTO(int64(len(todos)), offset, todos))

}

// @Summary Update a todo by ID
// @Description Update a todo by ID
// @Tags todo
// @Accept json
// @Produce json
// @Param id path string true "Todo ID"
// @Param todo body dtos.TodoDTO true "Todo object"
// @Success 204 {object} models.Todo
// @Router /todo/{id} [put]
func (t *todoHandler) Update(c *gin.Context) {

	id := c.Param("id")
	var todoDTO dtos.TodoDTO
	if err := c.ShouldBindJSON(&todoDTO); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	todo := todoDTO.ToModel()
	if err := t.todoDAO.Update(c, id, todo); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
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
// @Router /todo/{id} [delete]
func (t *todoHandler) Delete(c *gin.Context) {

	id := c.Param("id")
	if err := t.todoDAO.Delete(c, id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(204, nil)
}
