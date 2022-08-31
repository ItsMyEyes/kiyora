package routers

import (
	"myself_framwork/library/logger/v2"
	"myself_framwork/models"

	homeInit "myself_framwork/api/v1/controllers/home"

	"github.com/gin-gonic/gin"
)

func Routing(router *gin.Engine) *gin.Engine {
	gen := &models.GeneralModel{
		Logging: logger.InitLog(),
	}
	homeHandler := homeInit.InitiateHomeInterface(gen)

	api := router.Group("/api/")
	v1 := api.Group("/v1/")
	{
		home := v1.Group("/home")
		{
			home.GET("/", homeHandler.GetHome)
		}
	}

	return router
}
