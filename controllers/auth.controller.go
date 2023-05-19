package controllers

import (
	"errors"
	"net/http"

	"github.com/AmadoJunior/Gipitty/config"
	"github.com/AmadoJunior/Gipitty/models"
	"github.com/AmadoJunior/Gipitty/services"
	"github.com/AmadoJunior/Gipitty/utils"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
	userService services.UserService
}

func NewAuthController(authService services.AuthService, userService services.UserService) AuthController {
	return AuthController{authService, userService}
}

func (ac *AuthController) SignUpUser(ctx *gin.Context) {
	var user *models.SignUpInput

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if user.Password != user.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "passwords do not match"})
		return
	}

	newUser, err := ac.authService.SignUpUser(user)

	if err != nil {
		go utils.LogError(err, ctx)
		if errors.Is(err, services.ErrCreatingUser) {
			ctx.JSON(http.StatusConflict, gin.H{"status": "error", "message": "could not create your account"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "internal server error"})
		return
	}

	err = ac.userService.SendVerificationEmail(newUser)

	if err != nil {
		go utils.LogError(err, ctx)
		//Failed to Load Config
		if errors.Is(err, services.ErrLoadingConfig) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "internal server error"})
			return
		}

		//Failed to Send Email
		if errors.Is(err, services.ErrSendingEmail) {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "success", "message": "there was an error sending email"})
			return
		}

		//Failed to Update User Verification Code
		if errors.Is(err, services.ErrUpdateVerificationCode) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "internal server error"})
			return
		}

	}

	message := "we sent an email with a verification code to " + user.Email

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": message})
}

func (ac *AuthController) SignInUser(ctx *gin.Context) {
	var credentials *models.SignInInput

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	config, err := config.LoadConfig(".")

	if err != nil {
		go utils.LogError(err, ctx)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "success", "message": "internal server error"})
		return
	}

	access_token, refresh_token, err := ac.authService.SignInUser(credentials, config)

	if err != nil {
		go utils.LogError(err, ctx)
		//User Not Found || Incorrect Password
		if errors.Is(err, services.ErrUserNotFound) || errors.Is(err, services.ErrIncorrectPassword) {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "invalid email or password"})
			return
		}
		//Not Verified
		if errors.Is(err, services.ErrUserNotVerified) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "you are not verified, please verify your email to login"})
			return
		}

		//Failed to Generate Tokens
		if errors.Is(err, services.ErrGeneratingToken) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "success", "message": "internal server error"})
			return
		}

	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", config.Origin, false, true)
	ctx.SetCookie("refresh_token", refresh_token, config.RefreshTokenMaxAge*60, "/", config.Origin, false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", config.Origin, false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

func (ac *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	refresh_token, err := ctx.Cookie("refresh_token")

	if err != nil {
		go utils.LogError(err, ctx)
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	config, err := config.LoadConfig(".")

	if err != nil {
		go utils.LogError(err, ctx)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "success", "message": "internal server error"})
		return
	}

	access_token, err := ac.authService.RefreshAccessToken(refresh_token, config)

	if err != nil {
		go utils.LogError(err, ctx)
		//Invalid Token
		if errors.Is(err, services.ErrInvalidRefreshToken) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
			return
		}

		//User Not Found
		if errors.Is(err, services.ErrUserNotFound) {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
			return
		}

		//Failed to Create Token
		if errors.Is(err, services.ErrGeneratingToken) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "success", "message": "internal server error"})
			return
		}

	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", config.Origin, false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", config.Origin, false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

func (ac *AuthController) LogoutUser(ctx *gin.Context) {
	config, _ := config.LoadConfig(".")

	ctx.SetCookie("access_token", "", -1, "/", config.Origin, false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", config.Origin, false, true)
	ctx.SetCookie("logged_in", "", -1, "/", config.Origin, false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AuthController) VerifyEmail(ctx *gin.Context) {

	code := ctx.Params.ByName("verificationCode")

	err := ac.userService.VerifyUserEmail(code)
	if err != nil {
		go utils.LogError(err, ctx)
		ctx.JSON(http.StatusForbidden, gin.H{"status": "success", "message": "could not verify email address"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "email verified successfully"})

}

func (ac *AuthController) ForgotPassword(ctx *gin.Context) {
	var userCredential *models.ForgotPasswordInput

	if err := ctx.ShouldBindJSON(&userCredential); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "you will receive a reset email if user with that email exist"

	user, err := ac.userService.FindUserByEmail(userCredential.Email)
	if err != nil {
		go utils.LogError(err, ctx)
		if errors.Is(err, services.ErrUserEmailNotFound) {
			ctx.JSON(http.StatusOK, gin.H{"status": "fail", "message": message})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "success", "message": "internal server error"})
		return
	}

	if !user.Verified {
		ctx.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "account not verified"})
		return
	}

	config, err := config.LoadConfig(".")
	if err != nil {
		go utils.LogError(err, ctx)
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "success", "message": "internal server error"})
		return
	}

	// Update User in Database
	err = ac.userService.InitResetPassword(user, config)

	if err != nil {
		go utils.LogError(err, ctx)
		//User Not Found || Error Sending Mail
		if errors.Is(err, services.ErrUserEmailNotFound) || errors.Is(err, services.ErrSendingEmail) {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "success", "message": "there was an error sending email"})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

func (ac *AuthController) ResetPassword(ctx *gin.Context) {
	resetToken := ctx.Params.ByName("resetToken")
	var userCredential *models.ResetPasswordInput

	if err := ctx.ShouldBindJSON(&userCredential); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if userCredential.Password != userCredential.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "passwords do not match"})
		return
	}

	// Update User in Database
	err := ac.userService.ResetUserPassword(resetToken, userCredential.Password)

	if err != nil {
		go utils.LogError(err, ctx)
		if errors.Is(err, services.ErrResetTokenNotFound) {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "success", "message": "token is invalid or has expired"})
			return
		}
		if errors.Is(err, services.ErrUpdatingPassword) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"status": "success", "message": "internal server error"})
			return
		}
	}

	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "password data updated successfully"})
}
