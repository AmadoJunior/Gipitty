package repos

import "github.com/AmadoJunior/Gipitty/models"

type UserUpdate[T any] struct {
	Key   string
	Value T
}

type IUserRepo interface {
	//Core
	connect(dbUri string) error
	InitRepository(dbUri string, dbName string, repoName string) error
	DeinitRepository() error

	//Public
	CreateNewUser(*models.SignUpInput) (string, error)
	FindUserByID(id string) (*models.DBResponse, error)
	FindUserByEmail(email string) (*models.DBResponse, error)
	UpdateUserById(id string, update *models.UpdateInput) error
	UpdateUserByEmail(email string, update *models.UpdateInput) error
	VerifyUserEmail(verificationCode string) error
	StorePasswordResetToken(userEmail string, passwordResetToken string) error
	ResetUserPassword(passwordResetToken string, newPassword string) error
}
