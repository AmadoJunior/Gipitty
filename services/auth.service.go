package services

import (
	"github.com/AmadoJunior/Gipitty/config"
	"github.com/AmadoJunior/Gipitty/models"
)

type IAuthService interface {
	SignUpUser(*models.SignUpInput) (*models.DBResponse, error)
	SignInUser(*models.SignInInput, *config.Config) (string, string, error)
	RefreshAccessToken(string, *config.Config) (string, error)
}
