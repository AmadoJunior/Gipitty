package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AmadoJunior/Gipitty/models"
	"github.com/AmadoJunior/Gipitty/repos"
	"github.com/AmadoJunior/Gipitty/utils"
)

type AuthServiceImpl struct {
	UserRepo repos.IUserRepo
	ctx      context.Context
}

func NewAuthServiceImpl(userRepo repos.IUserRepo, ctx context.Context) AuthService {
	return &AuthServiceImpl{userRepo, ctx}
}

func (uc *AuthServiceImpl) SignUpUser(user *models.SignUpInput) (*models.DBResponse, error) {
	//Hash Password
	hashedPassword := make(chan string)
	errorChannel := make(chan error)
	go func(password string) {
		defer close(hashedPassword)
		defer close(errorChannel)
		result, err := utils.HashPassword(password)
		if err != nil {
			errorChannel <- err
		}
		hashedPassword <- result
	}(user.Password)

	//Init
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Email = strings.ToLower(user.Email)
	user.Verified = false
	user.Role = "user"
	user.PasswordConfirm = ""

	select {
	case err := <-errorChannel:
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrSignUp, err)
		}
	case user.Password = <-hashedPassword:
	}

	//Create User
	userId, err := uc.UserRepo.CreateNewUser(user)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrSignUp, err)
	}

	var newUser *models.DBResponse
	newUser, err = uc.UserRepo.FindUserByID(userId)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrSignUp, err)
	}

	return newUser, nil
}

func (uc *AuthServiceImpl) SignInUser(*models.SignInInput) (*models.DBResponse, error) {
	return nil, nil
}
