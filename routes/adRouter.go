package routes

import (
	controller "github.com/bunkieproject/bunkie_be/controllers"

	"github.com/gin-gonic/gin"
)

func AdRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("ads/create_bunkie", controller.CreateBunkieAd())
	incomingRoutes.GET("ads/get_bunkie", controller.GetBunkieAd())
	incomingRoutes.PUT("ads/update_bunkie", controller.UpdateBunkieAd())
	incomingRoutes.DELETE("ads/delete_bunkie", controller.DeleteBunkieAd())
	incomingRoutes.POST("ads/search_bunkie_default", controller.SearchBunkieAdDefault())
	incomingRoutes.POST("ads/search_bunkie_preferred", controller.SearchBunkieAdPreferred())
	incomingRoutes.POST("ads/create_room_ad", controller.CreateRoomAd())
	incomingRoutes.GET("ads/get_room_ad", controller.GetRoomAd())
	incomingRoutes.GET("ads/get_room_ads", controller.GetRoomAds())
	incomingRoutes.PUT("ads/update_room_ad", controller.UpdateRoomAd())
	incomingRoutes.DELETE("ads/delete_room_ad", controller.DeleteRoomAd())
	incomingRoutes.POST("ads/search_room_ad_default", controller.SearchRoomAdDefault())
	incomingRoutes.POST("ads/search_room_ad_preferred", controller.SearchRoomAdPreferred())
}
