package services

import (
	"github.com/AmadoJunior/Gipitty/models"
	"github.com/AmadoJunior/Gipitty/repos/userRepo"
)

type UserService interface {
	FindUserById(string) (*models.DBResponse, error)
	FindUserByEmail(string) (*models.DBResponse, error)
	UpdateUserById(id string, field string, value string) (*models.DBResponse, error)
	UpdateOne(field string, value interface{}) (*models.DBResponse, error)
	VerifyEmail(verificationCode string) (*userRepo.UpdatedResult, error)
	StorePasswordResetToken(userEmail string, passwordResetToken string) (*userRepo.UpdatedResult, error)
	ResetPassword(passwordResetToken string, newPassword string) (*userRepo.UpdatedResult, error)
}
