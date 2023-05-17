package services

import (
	"context"

	"github.com/AmadoJunior/Gipitty/models"
	"github.com/AmadoJunior/Gipitty/repos"
)

type UserServiceImpl struct {
	userRepo repos.IUserRepo
	ctx      context.Context
}

func NewUserServiceImpl(userRepo repos.IUserRepo, ctx context.Context) UserService {
	return &UserServiceImpl{userRepo, ctx}
}

func (us UserServiceImpl) FindUserById(id string) (*models.DBResponse, error) {
	return us.userRepo.FindUserByID(id)
}

func (us UserServiceImpl) FindUserByEmail(email string) (*models.DBResponse, error) {
	return us.userRepo.FindUserByEmail(email)
}

func (us UserServiceImpl) UpdateUserById(id string, data *models.UpdateInput) (*models.DBResponse, error) {
	err := us.userRepo.UpdateUserById(id, data)
	if err != nil {
		return nil, err
	}
	return us.userRepo.FindUserByID(id)
}

func (us UserServiceImpl) VerifyUserEmail(verificationCode string) error {
	return us.userRepo.VerifyUserEmail(verificationCode)
}

func (us UserServiceImpl) StorePasswordResetToken(userEmail string, passwordResetToken string) error {
	return us.userRepo.StorePasswordResetToken(userEmail, passwordResetToken)
}

func (us UserServiceImpl) ResetUserPassword(passwordResetToken string, newPassword string) error {
	return us.userRepo.ResetUserPassword(passwordResetToken, newPassword)
}
