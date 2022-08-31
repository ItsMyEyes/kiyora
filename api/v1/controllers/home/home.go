package home

import (
	"encoding/json"
	"myself_framwork/utils/validate"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (t *homeController) GetHome(ctx *gin.Context) {
	var query struct {
		ID    string `json"-"`
		Name  string `form:"name" json:"name" binding:"required,min=2"`
		Color string `form:"color" json:"color" binding:"required"`
	}

	if err := ctx.ShouldBind(&query); err != nil {
		t.Res.Code = http.StatusBadRequest
		t.Res.Message = validate.ParseError(err)
		t.Res.Status = false
		ctx.JSON(t.Res.Code, t.Res)
		return
	}

	res := t.homeService.Index()
	marsReq, _ := json.Marshal(query)
	marsRes, _ := json.Marshal(res)
	t.Gen.Logging.Info("API - Transaction (Request)", t.Gen.Logging.AddField("RequestBody", marsReq), t.Gen.Logging.AddField("ResponseBody", marsRes))
	ctx.JSON(200, res)
}
