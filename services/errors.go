package services

import "errors"

var (
	ErrSignUp                 = errors.New("failed to sign up user")
	ErrLoadingConfig          = errors.New("failed to load config")
	ErrUpdateVerificationCode = errors.New("failed to udpate user verification code")
	ErrSendingEmail           = errors.New("failed to send email")
)
