package database

import (
	"context"
	"todo-app-mongo/internal/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TodoDAOInterface interface {
	Create(ctx context.Context, todo *entity.Todo) error
	Get(ctx context.Context, id string, userId string) (*entity.Todo, error)
	GetAll(ctx context.Context, limit int64, page int64, search string, userId primitive.ObjectID) ([]*entity.Todo, int64, error)
	Update(ctx context.Context, id string, todo *entity.Todo) error
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

func (t *todoDAO) Create(ctx context.Context, todo *entity.Todo) error {
	_, err := t.collection.InsertOne(ctx, todo)
	return err
}

func (t *todoDAO) Get(ctx context.Context, id string, userId string) (*entity.Todo, error) {

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var todo *entity.Todo
	err = t.collection.FindOne(ctx, bson.M{"_id": objectID, "user_id": userId}).Decode(&todo)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (t *todoDAO) GetAll(ctx context.Context, limit int64, page int64, search string, userId primitive.ObjectID) ([]*entity.Todo, int64, error) {
	var todos []*entity.Todo

    filter := bson.M{"user_id": userId}
	if search != "" {
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": primitive.Regex{Pattern: search, Options: "i"}}},
			{"description": bson.M{"$regex": primitive.Regex{Pattern: search, Options: "i"}}},
		}
	}

	opts := options.Find()
	opts.SetLimit(limit)
	opts.SetSkip(page)
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := t.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var todo *entity.Todo
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

func (t *todoDAO) Update(ctx context.Context, id string, todo *entity.Todo) error {
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
