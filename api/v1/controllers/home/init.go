package home

import (
	"myself_framwork/api/v1/services/home"
	"myself_framwork/models"
)

type (
	homeController struct {
		Gen         *models.GeneralModel
		Res         models.Response
		homeService *home.HomeService
	}
)

func InitiateHomeInterface(gen *models.GeneralModel) *homeController {
	return &homeController{
		Gen:         gen,
		Res:         models.Response{},
		homeService: home.InitiateHomeServiceInterface(gen),
	}
}
