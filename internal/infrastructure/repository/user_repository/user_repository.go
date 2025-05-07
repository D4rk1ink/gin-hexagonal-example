package user_repository

import (
	"context"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/port"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/database"
	repository_mapper "github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/repository/mapper"
	repository_model "github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/repository/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type userRepository struct {
	mongodb database.MongoDb
}

func NewUserRepository(mongodb database.MongoDb) port.UserRepository {
	return &userRepository{
		mongodb: mongodb,
	}
}

func (r *userRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	cursor, err := r.mongodb.GetDb().Collection("users").Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*domain.User

	for cursor.Next(ctx) {
		var user repository_model.UserModel
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, repository_mapper.ToUserDomain(&user))
	}

	return users, nil
}

func (r *userRepository) GetById(ctx context.Context, id string) (*domain.User, error) {
	var userModel repository_model.UserModel

	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = r.mongodb.GetDb().Collection("users").FindOne(ctx, bson.M{"_id": objID}).Decode(&userModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return repository_mapper.ToUserDomain(&userModel), nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var userModel repository_model.UserModel

	err := r.mongodb.GetDb().Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&userModel)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return repository_mapper.ToUserDomain(&userModel), nil
}

func (r *userRepository) Create(ctx context.Context, payload *domain.User) (*string, error) {
	result, err := r.mongodb.GetDb().Collection("users").InsertOne(ctx, repository_mapper.ToUserModel(payload))
	if err != nil {
		return nil, err
	}

	insertedId := result.InsertedID.(bson.ObjectID).Hex()

	return &insertedId, nil
}

func (r *userRepository) Update(ctx context.Context, payload *domain.User) error {
	_, err := r.mongodb.GetDb().Collection("users").UpdateByID(ctx, payload.ID, repository_mapper.ToUserModel(payload))
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return nil
}
