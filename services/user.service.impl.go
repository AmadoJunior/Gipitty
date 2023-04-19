package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AmadoJunior/Gipitty/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserServiceImpl(collection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{collection, ctx}
}

func (us *UserServiceImpl) FindUserById(id string) (*models.DBResponse, error) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var user *models.DBResponse

	query := bson.M{"_id": oid}
	err := us.collection.FindOne(us.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) FindUserByEmail(email string) (*models.DBResponse, error) {
	var user *models.DBResponse

	query := bson.M{"email": strings.ToLower(email)}
	err := us.collection.FindOne(us.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

func (uc *UserServiceImpl) UpdateUserById(id string, field string, value string) (*models.DBResponse, error) {
	userId, idErr := primitive.ObjectIDFromHex(id)
	if idErr != nil {
		fmt.Print("ID ERR", idErr)
		return &models.DBResponse{}, idErr
	}
	query := bson.D{{Key: "_id", Value: userId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: value}}}}
	result, err := uc.collection.UpdateOne(uc.ctx, query, update)

	if err != nil {
		fmt.Print(err)
		return &models.DBResponse{}, err
	}
	fmt.Print(result.ModifiedCount)
	return &models.DBResponse{}, nil
}

func (uc *UserServiceImpl) UpdateOne(field string, value interface{}) (*models.DBResponse, error) {
	query := bson.D{{Key: field, Value: value}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: value}}}}
	result, err := uc.collection.UpdateOne(uc.ctx, query, update)

	if err != nil {
		fmt.Print(err)
		return &models.DBResponse{}, err
	}
	fmt.Print(result.ModifiedCount)
	return &models.DBResponse{}, nil
}

func (uc *UserServiceImpl) VerifyEmail(verificationCode string) (*mongo.UpdateResult, error) {
	query := bson.D{{Key: "verificationCode", Value: verificationCode}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "verified", Value: true}}}, {Key: "$unset", Value: bson.D{{Key: "verificationCode", Value: ""}}}}
	result, err := uc.collection.UpdateOne(uc.ctx, query, update)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	fmt.Print(result.ModifiedCount)
	return result, nil
}

func (uc *UserServiceImpl) StorePasswordResetToken(userEmail string, passwordResetToken string) (*mongo.UpdateResult, error) {
	query := bson.D{{Key: "email", Value: strings.ToLower(userEmail)}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "passwordResetToken", Value: passwordResetToken}, {Key: "passwordResetAt", Value: time.Now().Add(time.Minute * 15)}}}}
	result, err := uc.collection.UpdateOne(uc.ctx, query, update)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	fmt.Print(result.ModifiedCount)
	return result, nil
}

func (uc *UserServiceImpl) ResetPassword(passwordResetToken string, newPassword string) (*mongo.UpdateResult, error) {
	query := bson.D{{Key: "passwordResetToken", Value: passwordResetToken}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: passwordResetToken}}}, {Key: "$unset", Value: bson.D{{Key: "passwordResetToken", Value: ""}, {Key: "passwordResetAt", Value: ""}}}}
	result, err := uc.collection.UpdateOne(uc.ctx, query, update)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	fmt.Print(result.ModifiedCount)
	return result, nil
}
