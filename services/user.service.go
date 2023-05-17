package services

import (
	"github.com/AmadoJunior/Gipitty/models"
	"github.com/AmadoJunior/Gipitty/repos"
)

type UserService interface {
	FindUserById(string) (*models.DBResponse, error)
	FindUserByEmail(string) (*models.DBResponse, error)
	UpdateUserById(id string, field string, value string) (*models.DBResponse, error)
	UpdateOne(field string, value interface{}) (*models.DBResponse, error)
	VerifyEmail(verificationCode string) (*repos.UpdatedResult, error)
	StorePasswordResetToken(userEmail string, passwordResetToken string) (*repos.UpdatedResult, error)
	ResetPassword(passwordResetToken string, newPassword string) (*repos.UpdatedResult, error)
}
