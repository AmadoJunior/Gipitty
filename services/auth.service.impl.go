package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AmadoJunior/Gipitty/config"
	"github.com/AmadoJunior/Gipitty/models"
	"github.com/AmadoJunior/Gipitty/repos"
	"github.com/AmadoJunior/Gipitty/utils"
)

type AuthService struct {
	UserRepo repos.IUserRepo
	ctx      context.Context
}

func NewAuthService(userRepo repos.IUserRepo, ctx context.Context) IAuthService {
	return &AuthService{userRepo, ctx}
}

func (uc *AuthService) SignUpUser(user *models.SignUpInput) (*models.DBResponse, error) {
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
			return nil, utils.GenerateError(ErrHashingPassword, err)
		}
	case user.Password = <-hashedPassword:
	}

	//Create User
	userId, err := uc.UserRepo.CreateNewUser(user)
	if err != nil {
		return nil, utils.GenerateError(ErrCreatingUser, err)
	}

	var newUser *models.DBResponse
	newUser, err = uc.UserRepo.FindUserByID(userId)
	if err != nil {
		return nil, utils.GenerateError(ErrUserIDNotFound, err)
	}

	return newUser, nil
}

func (uc *AuthService) SignInUser(credentials *models.SignInInput, config *config.Config) (string, string, error) {
	user, err := uc.UserRepo.FindUserByEmail(credentials.Email)
	if err != nil {
		//User Not Found
		return "", "", utils.GenerateError(ErrUserNotFound, err)
	}

	if !user.Verified {
		//Not Verified
		return "", "", utils.GenerateError(ErrUserNotVerified, err)
	}

	if err := utils.VerifyPassword(user.Password, credentials.Password); err != nil {
		//Incorrect Password
		return "", "", utils.GenerateError(ErrIncorrectPassword, err)
	}

	// Generate Tokens
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		//Failed to Generate Tokens
		return "", "", utils.GenerateError(ErrGeneratingToken, err)
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		//Failed to Generate Tokens
		return "", "", utils.GenerateError(ErrGeneratingToken, err)
	}

	return access_token, refresh_token, nil
}

func (uc AuthService) RefreshAccessToken(refresh_token string, config *config.Config) (string, error) {
	sub, err := utils.ValidateToken(refresh_token, config.RefreshTokenPublicKey)
	if err != nil {
		//Invalid Token
		return "", utils.GenerateError(ErrInvalidRefreshToken, err)
	}

	user, err := uc.UserRepo.FindUserByID(fmt.Sprint(sub))
	if err != nil {
		//User Not Found
		return "", utils.GenerateError(ErrUserNotFound, err)
	}

	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		//Failed to Create Token
		return "", utils.GenerateError(ErrGeneratingToken, err)
	}

	return access_token, nil
}
