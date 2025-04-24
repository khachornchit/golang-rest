package mongo_repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang-rest/internal/core/domain"
	"golang-rest/internal/core/ports"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) ports.UserRepositoryInterface {
	return &UserRepository{collection: collection}
}

func (u UserRepository) CreateUser(user *domain.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Password hashing failed: ", err)
		return err
	}
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()

	_, err = u.collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Println("MongoDB insert error: ", err)
	}

	return err
}

func (u UserRepository) GetAllUsers() ([]domain.User, error) {
	projection := bson.D{{Key: "password", Value: 0}}
	findOptions := options.Find().SetProjection(projection)
	cursor, err := u.collection.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var users []domain.User
	if err := cursor.All(context.Background(), &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (u UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	projection := bson.D{{Key: "password", Value: 0}}
	findOptions := options.FindOne().SetProjection(projection)
	err := u.collection.FindOne(context.Background(), bson.M{"email": email}, findOptions).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u UserRepository) GetUserLoginByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := u.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u UserRepository) GetUserByID(id string) (*domain.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user domain.User
	projection := bson.D{{Key: "password", Value: 0}}
	findOptions := options.FindOne().SetProjection(projection)
	err = u.collection.FindOne(context.Background(), bson.M{"_id": objectID}, findOptions).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u UserRepository) UpdateUserByID(id string, updates bson.M) (*domain.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	update := bson.M{"$set": updates}
	_, err = u.collection.UpdateByID(context.Background(), objectID, update)
	if err != nil {
		return nil, err
	}

	return u.GetUserByID(id)
}

func (u UserRepository) DeleteUserByID(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = u.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}
