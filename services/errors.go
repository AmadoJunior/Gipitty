package services

import "errors"

var (
	ErrCreatingUser           = errors.New("failed to create user")
	ErrLoadingConfig          = errors.New("failed to load config")
	ErrUpdateVerificationCode = errors.New("failed to udpate user verification code")
	ErrSendingEmail           = errors.New("failed to send email")
	ErrUserNotFound           = errors.New("failed to find user")
	ErrUserNotVerified        = errors.New("user not verified")
	ErrIncorrectPassword      = errors.New("incorrect password")
	ErrGeneratingToken        = errors.New("failed to generate authentication tokens")
	ErrInvalidRefreshToken    = errors.New("invalid refresh token")
	ErrUserEmailNotFound      = errors.New("failed to find user by email")
	ErrUserIDNotFound         = errors.New("failed to find user by id")
	ErrHashingPassword        = errors.New("failed to hash password")
	ErrResetTokenNotFound     = errors.New("failed to find user with reset token")
	ErrUpdatingPassword       = errors.New("failed to update password")
)
