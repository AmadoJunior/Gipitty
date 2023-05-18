package services

import (
	"github.com/AmadoJunior/Gipitty/config"
	"github.com/AmadoJunior/Gipitty/models"
)

type UserService interface {
	FindUserById(id string) (*models.DBResponse, error)
	FindUserByEmail(email string) (*models.DBResponse, error)
	UpdateUserById(id string, data *models.UpdateInput) (*models.DBResponse, error)
	SendVerificationEmail(newUser *models.DBResponse) error
	VerifyUserEmail(verificationCode string) error
	InitResetPassword(*models.DBResponse, config.Config) error
	ResetUserPassword(passwordResetToken string, newPassword string) error
}
