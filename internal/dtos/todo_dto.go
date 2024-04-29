package dtos

import (
	"time"
	"todo-app-mongo/internal/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TodoDTO struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Scheduled   bool   `json:"scheduled"`
	ScheduledTo string `json:"scheduled_to"`
}

func (t *TodoDTO) ToModel() *entity.Todo {
	model := &entity.Todo{
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

func (t *TodoDTO) FromModel(todo *entity.Todo) {
	t.Title = todo.Title
	t.Description = todo.Description
	t.Scheduled = todo.Scheduled

	if !todo.ScheduledTo.IsZero() {
		t.ScheduledTo = todo.ScheduledTo.Format(time.RFC3339)
	}
}

func (t *TodoDTO) ToModelUpdate() *entity.Todo {
	model := &entity.Todo{
		Title:       t.Title,
		Description: t.Description,
		Scheduled:   t.Scheduled,
	}

	if t.ScheduledTo != "" {
		scheduledTo, _ := time.Parse(time.RFC3339, t.ScheduledTo)
		model.ScheduledTo = scheduledTo
	}

	return model
}
