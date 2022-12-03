package routes

import (
	controller "bunkie_be/controllers"

	"github.com/gin-gonic/gin"
)

func AdRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("ads/create", controller.CreateHouseAd())
	incomingRoutes.POST("ads/delete", controller.DeleteHouseAd())
	incomingRoutes.POST("ads/update", controller.UpdateHouseAd())
	incomingRoutes.POST("ads/create_bunkie", controller.CreateBunkieAd())
	incomingRoutes.POST("ads/delete_bunkie", controller.DeleteBunkieAd())
	incomingRoutes.POST("ads/update_bunkie", controller.UpdateBunkieAd())
}
