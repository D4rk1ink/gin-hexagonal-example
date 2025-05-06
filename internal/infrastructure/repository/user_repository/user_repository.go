package user_repository

import (
	"context"

	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/domain"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/dto"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/core/port"
	"github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/database"
	repository_mapper "github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/repository/mapper"
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
	cursor, err := r.mongodb.GetDb().Collection("users").Find(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []*domain.User

	for cursor.Next(ctx) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *userRepository) GetById(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User

	err := r.mongodb.GetDb().Collection("users").FindOne(ctx, map[string]interface{}{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	err := r.mongodb.GetDb().Collection("users").FindOne(ctx, map[string]interface{}{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, payload *domain.User) (*string, error) {
	result, err := r.mongodb.GetDb().Collection("users").InsertOne(ctx, repository_mapper.ToUserModel(payload))

	if err != nil {
		return nil, err
	}

	insertedId := result.InsertedID.(bson.ObjectID).Hex()

	return &insertedId, nil
}

func (r *userRepository) Update(ctx context.Context, payload dto.UserUpdateDto) error {
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	return nil
}
