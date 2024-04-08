package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Todo struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Scheduled   bool               `json:"scheduled" bson:"scheduled"`
	ScheduledTo time.Time          `json:"scheduled_to" bson:"scheduled_to"`
	Completed   bool               `json:"completed" bson:"completed"`
	CompletedAt time.Time          `json:"completed_at" bson:"completed_at"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UserID      primitive.ObjectID `json:"user_id" bson:"user_id"`
}
