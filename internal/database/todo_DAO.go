package database

import (
	"context"
	"todo-app-mongo/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TodoDAOInterface interface {
	Create(ctx context.Context, todo *models.Todo) error
	Get(ctx context.Context, id string) (*models.Todo, error)
	GetAll(ctx context.Context, limit int64, page int64, search string) ([]*models.Todo, int64, error)
	Update(ctx context.Context, id string, todo *models.Todo) error
	Delete(ctx context.Context, id string) error
}

type todoDAO struct {
	collection *mongo.Collection
}

func NewTodoDAO(db mongo.Database) *todoDAO {
	return &todoDAO{
		collection: db.Collection("todos"),
	}
}

func (t *todoDAO) Create(ctx context.Context, todo *models.Todo) error {
	_, err := t.collection.InsertOne(ctx, todo)
	return err
}

func (t *todoDAO) Get(ctx context.Context, id string) (*models.Todo, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var todo *models.Todo
	err = t.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&todo)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (t *todoDAO) GetAll(ctx context.Context, limit int64, page int64, search string) ([]*models.Todo, int64, error) {
	var todos []*models.Todo

	filter := bson.M{}
	if search != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": primitive.Regex{Pattern: search, Options: "i"}}},
			{"description": bson.M{"$regex": primitive.Regex{Pattern: search, Options: "i"}}},
		}
	}

	cursor, err := t.collection.Find(ctx, filter, &options.FindOptions{Limit: &limit, Skip: &page})
	if err != nil {
		return nil, 0, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var todo *models.Todo
		if err := cursor.Decode(&todo); err != nil {
			return nil, 0, err
		}
		todos = append(todos, todo)
	}

	count, err := t.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	//calculate total pages
	totalPages := (count / limit)

	return todos, totalPages, nil
}

func (t *todoDAO) Update(ctx context.Context, id string, todo *models.Todo) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = t.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": todo})
	return err
}

func (t *todoDAO) Delete(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = t.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
