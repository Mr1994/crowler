package v1

import (
	"api_client/help"
	"api_client/request"
	"api_client/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type HomeController struct {
}

func (p *HomeController) HomeInfo(c *gin.Context) {
	// 获取请求参数
	var params request.HomeParams

	_ = c.ShouldBindBodyWith(&params, binding.JSON)
	var homeSelectParams request.HomeSelectParams
	homeSelectParams.PageSize = 10
	homeSelectParams.NativePlace = params.NativePlace

	foundService := new(service.FoundPersonService)
	foundPetService := new(service.FoundPetService)
	personList, err := foundService.GetHomePersonList(&homeSelectParams)
	if err != nil {
		panic("系统异常:" + err.Error())
	}
	petList, err := foundPetService.GetHomePetList(&homeSelectParams)
	if err != nil {
		panic("系统异常:" + err.Error())
	}
	var listMap map[string]interface{}
	listMap = make(map[string]interface{})
	listMap["personList"] = personList
	listMap["petList"] = petList
	// 首页开关
	listMap["optionSwitch"] = map[string]bool{
		"foundPerson": true,
		"foundPet":    true,
		"foundThing":  false,
		"adopt":       false,
	}
	listMap["bannerImage"] = []string{
		"http://image.foundall.cn/static/gonggao.png",
		"http://image.foundall.cn/static/gonggao2.png",
	}
	data := map[string]interface{}{
		"list": listMap,
	}
	help.Output(data, c)
	return

}
