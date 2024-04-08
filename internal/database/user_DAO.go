package database

import (
	"context"
	"time"
	"todo-app-mongo/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserDAOInterface interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, email string) error
	GetById(ctx context.Context, id string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
}

type userDAO struct {
	collection *mongo.Collection
}

func NewUserDAO(db mongo.Database) *userDAO {
	return &userDAO{
		collection: db.Collection("todo_user"),
	}

}

func (u *userDAO) Create(ctx context.Context, user *models.User) error {

	_, err := u.collection.InsertOne(ctx, user)
	return err
}

func (u *userDAO) Update(ctx context.Context, user *models.User) error {

	_, err := u.collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
	return err
}

func (u *userDAO) Delete(ctx context.Context, email string) error {

	user, err := u.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	user.Removed = true
	user.RemovedAt = time.Now()
	user.UpdatedAt = time.Now()

	return u.Update(ctx, user)

}

func (u *userDAO) GetById(ctx context.Context, id string) (*models.User, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user *models.User
	err = u.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userDAO) GetByEmail(ctx context.Context, email string) (*models.User, error) {

	var user *models.User
	err := u.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
