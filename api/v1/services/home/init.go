package home

import (
	"myself_framwork/models"
)

type (
	HomeService struct {
		Gen *models.GeneralModel
	}
)

func InitiateHomeServiceInterface(gen *models.GeneralModel) *HomeService {
	return &HomeService{
		Gen: gen,
	}
}
