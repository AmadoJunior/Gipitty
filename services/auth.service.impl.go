package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/AmadoJunior/Gipitty/models"
	"github.com/AmadoJunior/Gipitty/repos/userRepo"
	"github.com/AmadoJunior/Gipitty/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type AuthServiceImpl struct {
	userRepo userRepo.IUserRepo
	ctx      context.Context
}

func NewAuthServiceImpl(userRepo userRepo.IUserRepo, ctx context.Context) AuthService {
	return &AuthServiceImpl{userRepo, ctx}
}

func (uc *AuthServiceImpl) SignUpUser(user *models.SignUpInput) (*models.DBResponse, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Email = strings.ToLower(user.Email)
	user.PasswordConfirm = ""
	user.Verified = false
	user.Role = "user"

	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword
	res, err := uc.userRepo.InsertUser(&user)

	if err != nil {
		return nil, err
	}

	if key, err := uc.userRepo.CreateUserIndex("email", true); err != nil {
		return nil, errors.New("could not create index for user " + key)
	}

	var newUser *models.DBResponse = &models.DBResponse{}
	query := bson.M{"_id": res.InsertedUserID}
	err = uc.userRepo.FindUser(newUser, query)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (uc *AuthServiceImpl) SignInUser(*models.SignInInput) (*models.DBResponse, error) {
	return nil, nil
}
