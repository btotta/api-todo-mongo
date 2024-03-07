package dtos

import (
	"time"
	"todo-app-mongo/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoDTO struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Scheduled   bool   `json:"scheduled"`
	ScheduledTo string `json:"scheduled_to"`
}

func (t *TodoDTO) ToModel() *models.Todo {
	model := &models.Todo{
		ID:          primitive.NewObjectID(),
		Title:       t.Title,
		Description: t.Description,
		Scheduled:   t.Scheduled,
		CreatedAt:   time.Now(),
	}

	if t.ScheduledTo != "" {
		scheduledTo, _ := time.Parse(time.RFC3339, t.ScheduledTo)
		model.ScheduledTo = scheduledTo
	}

	return model
}

func (t *TodoDTO) FromModel(todo *models.Todo) {
	t.Title = todo.Title
	t.Description = todo.Description
	t.Scheduled = todo.Scheduled

	if !todo.ScheduledTo.IsZero() {
		t.ScheduledTo = todo.ScheduledTo.Format(time.RFC3339)
	}
}
