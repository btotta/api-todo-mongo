package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	Name           string             `json:"name" bson:"name"`
	Email          string             `json:"email" bson:"email"`
	HashedPassword string             `json:"hashed_password" bson:"hashed_password"`
	CreatedAt      time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at" bson:"updated_at"`
	Removed        bool               `json:"removed" bson:"removed"`
	RemovedAt      time.Time          `json:"removed_at" bson:"removed_at"`
}

func (u *User) ComparePassword(password string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password))
	return err == nil
}
