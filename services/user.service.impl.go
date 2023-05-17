package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AmadoJunior/Gipitty/models"
	"github.com/AmadoJunior/Gipitty/repos"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserServiceImpl struct {
	userRepo repos.IUserRepo
	ctx      context.Context
}

func NewUserServiceImpl(userRepo repos.IUserRepo, ctx context.Context) UserService {
	return &UserServiceImpl{userRepo, ctx}
}

func (us *UserServiceImpl) FindUserById(id string) (*models.DBResponse, error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	var user *models.DBResponse = &models.DBResponse{}
	query := bson.M{"_id": oid}
	err := us.userRepo.FindUser(user, query)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) FindUserByEmail(email string) (*models.DBResponse, error) {
	var user *models.DBResponse = &models.DBResponse{}

	query := bson.M{"email": strings.ToLower(email)}
	err := us.userRepo.FindUser(user, query)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserServiceImpl) UpdateUserById(id string, field string, value string) (*models.DBResponse, error) {
	userId, idErr := primitive.ObjectIDFromHex(id)
	if idErr != nil {
		return &models.DBResponse{}, idErr
	}
	query := bson.D{{Key: "_id", Value: userId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: value}}}}
	result, err := uc.userRepo.UpdateUser(query, update, false)

	if err != nil {
		return &models.DBResponse{}, err
	}
	fmt.Print(result.ModifiedCount)
	return &models.DBResponse{}, nil
}

func (uc *UserServiceImpl) UpdateOne(field string, value interface{}) (*models.DBResponse, error) {
	query := bson.D{{Key: field, Value: value}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: value}}}}
	result, err := uc.userRepo.UpdateUser(query, update, false)

	if err != nil {
		return &models.DBResponse{}, err
	}
	fmt.Print(result.ModifiedCount)
	return &models.DBResponse{}, nil
}

func (uc *UserServiceImpl) VerifyEmail(verificationCode string) (*repos.UpdatedResult, error) {
	query := bson.D{{Key: "verificationCode", Value: verificationCode}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "verified", Value: true}}}, {Key: "$unset", Value: bson.D{{Key: "verificationCode", Value: ""}}}}
	result, err := uc.userRepo.UpdateUser(query, update, false)

	if err != nil {
		return nil, err
	}
	fmt.Print(result.ModifiedCount)
	return result, nil
}

func (uc *UserServiceImpl) StorePasswordResetToken(userEmail string, passwordResetToken string) (*repos.UpdatedResult, error) {
	query := bson.D{{Key: "email", Value: strings.ToLower(userEmail)}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "passwordResetToken", Value: passwordResetToken}, {Key: "passwordResetAt", Value: time.Now().Add(time.Minute * 15)}}}}
	result, err := uc.userRepo.UpdateUser(query, update, false)

	if err != nil {
		return nil, err
	}
	fmt.Print(result.ModifiedCount)
	return result, nil
}

func (uc *UserServiceImpl) ResetPassword(passwordResetToken string, newPassword string) (*repos.UpdatedResult, error) {
	query := bson.D{{Key: "passwordResetToken", Value: passwordResetToken}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: passwordResetToken}}}, {Key: "$unset", Value: bson.D{{Key: "passwordResetToken", Value: ""}, {Key: "passwordResetAt", Value: ""}}}}
	result, err := uc.userRepo.UpdateUser(query, update, false)

	if err != nil {
		return nil, err
	}
	fmt.Print(result.ModifiedCount)
	return result, nil
}
