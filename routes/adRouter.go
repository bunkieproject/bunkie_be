package routes

import (
	controller "bunkie_be/controllers"

	"github.com/gin-gonic/gin"
)

func AdRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("ads/create_bunkie", controller.CreateBunkieAd())
	incomingRoutes.GET("ads/get_bunkie", controller.GetBunkieAd())
	incomingRoutes.PUT("ads/update_bunkie", controller.UpdateBunkieAd())
	incomingRoutes.DELETE("ads/delete_bunkie", controller.DeleteBunkieAd())
	incomingRoutes.GET("ads/search_bunkie_default", controller.SearchBunkieAdDefault())
	incomingRoutes.GET("ads/search_bunkie_preferred", controller.SearchBunkieAdPreferred())
	incomingRoutes.POST("ads/create_room_ad", controller.CreateRoomAd())
	incomingRoutes.GET("ads/get_room_ad", controller.GetRoomAd())
	incomingRoutes.PUT("ads/update_room_ad", controller.UpdateRoomAd())
	incomingRoutes.DELETE("ads/delete_room_ad", controller.DeleteRoomAd())
	incomingRoutes.GET("ads/search_room_ad_default", controller.SearchRoomAdDefault())
	incomingRoutes.GET("ads/search_room_ad_preferred", controller.SearchRoomAdPreferred())
}
