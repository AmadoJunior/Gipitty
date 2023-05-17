package repos

import "errors"

var (
	ErrUserRepoInit            = errors.New("failed to initiate user repository")
	ErrUserRepoDeinit          = errors.New("failed to deinitiate user repository")
	ErrDuplicateEmail          = errors.New("user email already in use")
	ErrUserInsertion           = errors.New("failed to insert user")
	ErrUserIDAssertion         = errors.New("failed to assert user object id")
	ErrInvalidIDHex            = errors.New("failed to create object id from hex")
	ErrUserNotFound            = errors.New("failed to find user")
	ErrUserUpdate              = errors.New("failed to update user")
	ErrUserVerification        = errors.New("failed to verify user")
	ErrStorePasswordResetToken = errors.New("failed to store password reset token")
	ErrResetPassword           = errors.New("failed to reset password")
	ErrInvalidUpdateInput      = errors.New("provided update input is invalid")
)
