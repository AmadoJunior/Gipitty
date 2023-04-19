package services

import (
	"github.com/AmadoJunior/Gipitty/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService interface {
	FindUserById(string) (*models.DBResponse, error)
	FindUserByEmail(string) (*models.DBResponse, error)
	UpdateUserById(id string, field string, value string) (*models.DBResponse, error)
	UpdateOne(field string, value interface{}) (*models.DBResponse, error)
	VerifyEmail(verificationCode string) (*mongo.UpdateResult, error)
	StorePasswordResetToken(userEmail string, passwordResetToken string) (*mongo.UpdateResult, error)
	ResetPassword(passwordResetToken string, newPassword string) (*mongo.UpdateResult, error)
}
