package routes

import (
	"github.com/AmadoJunior/Gipitty/controllers"
	"github.com/AmadoJunior/Gipitty/middleware"
	"github.com/AmadoJunior/Gipitty/services"
	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controllers.UserController
	userService    services.UserService
}

func NewRouteUserController(userController controllers.UserController, userService services.UserService) UserRouteController {
	return UserRouteController{userController, userService}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {

	router := rg.Group("/users")
	router.Use(middleware.DeserializeUser(uc.userService))
	router.GET("/me", uc.userController.GetMe)
}
