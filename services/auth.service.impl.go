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
			return nil, fmt.Errorf("%w: %w", ErrHashingPassword, err)
		}
	case user.Password = <-hashedPassword:
	}

	//Create User
	userId, err := uc.UserRepo.CreateNewUser(user)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrCreatingUser, err)
	}

	var newUser *models.DBResponse
	newUser, err = uc.UserRepo.FindUserByID(userId)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUserIDNotFound, err)
	}

	return newUser, nil
}

func (uc *AuthServiceImpl) SignInUser(credentials *models.SignInInput, config config.Config) (string, string, error) {
	user, err := uc.UserRepo.FindUserByEmail(credentials.Email)
	if err != nil {
		//User Not Found
		return "", "", fmt.Errorf("%w: %w", ErrUserNotFound, err)
	}

	if !user.Verified {
		//Not Verified
		return "", "", fmt.Errorf("%w: %w", ErrUserNotVerified, err)
	}

	if err := utils.VerifyPassword(user.Password, credentials.Password); err != nil {
		//Incorrect Password
		return "", "", fmt.Errorf("%w: %w", ErrIncorrectPassword, err)
	}

	// Generate Tokens
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		//Failed to Generate Tokens
		return "", "", fmt.Errorf("%w: %w", ErrGeneratingToken, err)
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		//Failed to Generate Tokens
		return "", "", fmt.Errorf("%w: %w", ErrGeneratingToken, err)
	}

	return access_token, refresh_token, nil
}

func (uc AuthServiceImpl) RefreshAccessToken(refresh_token string, config config.Config) (string, error) {
	sub, err := utils.ValidateToken(refresh_token, config.RefreshTokenPublicKey)
	if err != nil {
		//Invalid Token
		return "", fmt.Errorf("%w:%w", ErrInvalidRefreshToken, err)
	}

	user, err := uc.UserRepo.FindUserByID(fmt.Sprint(sub))
	if err != nil {
		//User Not Found
		return "", fmt.Errorf("%w:%w", ErrUserNotFound, err)
	}

	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		//Failed to Create Token
		return "", fmt.Errorf("%w:%w", ErrGeneratingToken, err)
	}

	return access_token, nil
}
