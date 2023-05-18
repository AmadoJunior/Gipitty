package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/AmadoJunior/Gipitty/config"
	"github.com/AmadoJunior/Gipitty/models"
	"github.com/AmadoJunior/Gipitty/repos"
	"github.com/AmadoJunior/Gipitty/utils"
	"github.com/thanhpk/randstr"
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

func (us UserServiceImpl) SendVerificationEmail(newUser *models.DBResponse) error {
	config, err := config.LoadConfig(".")
	if err != nil {
		//Error Loading Config
		return fmt.Errorf("%w: %w", ErrLoadingConfig, err)
	}

	// Generate Verification Code
	code := randstr.String(20)
	verificationCode := utils.Encode(code)

	//Update User Async
	errorChan := make(chan error)

	go func() {
		defer close(errorChan)
		// Update User in Database
		updateData := &models.UpdateInput{
			VerificationCode: verificationCode,
		}
		_, err = us.UpdateUserById(newUser.ID.Hex(), updateData)
		if err != nil {
			errorChan <- err
		}
		errorChan <- nil
	}()

	//Get firstName
	var firstName = newUser.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// Send Email
	emailData := utils.EmailData{
		URL:       config.Origin + "/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	err = utils.SendEmail(newUser, &emailData, "verificationCode.html")
	if err != nil {
		//Error Sending Mail
		return fmt.Errorf("%w: %w", ErrSendingEmail, err)
	}

	err = <-errorChan
	if err != nil {
		//Error Setting Verification Code
		return fmt.Errorf("%w: %w", ErrUpdateVerificationCode, err)
	}

	return nil
}

func (us UserServiceImpl) VerifyUserEmail(code string) error {
	verificationCode := utils.Encode(code)
	return us.userRepo.VerifyUserEmail(verificationCode)
}

func (us UserServiceImpl) StorePasswordResetToken(userEmail string) (string, error) {
	// Generate Verification Code
	resetToken := randstr.String(20)

	passwordResetToken := utils.Encode(resetToken)

	err := us.userRepo.StorePasswordResetToken(userEmail, passwordResetToken)

	if err != nil {
		return "", err
	}

	return resetToken, nil
}

func (us UserServiceImpl) ResetUserPassword(resetToken string, newPassword string) error {
	errChan := make(chan error)
	outChan := make(chan string)
	go func() {
		defer close(errChan)
		defer close(outChan)
		result, err := utils.HashPassword(newPassword)
		if err != nil {
			errChan <- err
		} else {
			outChan <- result
		}
	}()

	passwordResetToken := utils.Encode(resetToken)
	var hashedPassword string
	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
	case hashedPassword = <-outChan:
	}
	return us.userRepo.ResetUserPassword(passwordResetToken, hashedPassword)
}
