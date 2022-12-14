package routes

import (
	controller "bunkie_be/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controller.SignUp())
	incomingRoutes.POST("users/login", controller.Login())
	incomingRoutes.POST("users/logout", controller.Logout())
	incomingRoutes.POST("users/reset_password", controller.ResetPassword())
	incomingRoutes.POST("users/enter_new_password", controller.EnterNewPassword())
	incomingRoutes.DELETE("users/delete_account", controller.DeleteAccount())
	incomingRoutes.GET("users/get_profile_info", controller.GetProfileInfo())
	incomingRoutes.PUT("users/settings", controller.UpdateAccountInfo())
	incomingRoutes.POST("users/create_profile", controller.CreateProfileInfo())
	incomingRoutes.PUT("users/edit_profile", controller.EditProfileInfo())
	incomingRoutes.GET("users/display_profile", controller.DisplayProfile())
}
