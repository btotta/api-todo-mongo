package dtos

import (
	"time"
	"todo-app-mongo/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoDTO struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Scheduled   bool      `json:"scheduled"`
	ScheduledTo time.Time `json:"scheduled_to"`
}

func (t *TodoDTO) ToModel() *models.Todo {
	return &models.Todo{
		ID:          primitive.NewObjectID(),
		Title:       t.Title,
		Description: t.Description,
		Scheduled:   t.Scheduled,
		ScheduledTo: t.ScheduledTo,
		CreatedAt:   time.Now(),
	}
}

func (t *TodoDTO) FromModel(todo *models.Todo) {
	t.Title = todo.Title
	t.Description = todo.Description
	t.Scheduled = todo.Scheduled
	t.ScheduledTo = todo.ScheduledTo
}
