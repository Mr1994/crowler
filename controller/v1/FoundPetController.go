package v1

import (
	"api_client/help"
	"api_client/request"
	"api_client/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type FoundPetController struct {
}

/**
获取失踪人口用户列表
*/
func (p *FoundPetController) List(c *gin.Context) {
	// 获取请求参数
	var params *request.HomeSelectParams

	_ = c.ShouldBindBodyWith(&params, binding.JSON)
	if params.PageSize == 0 {
		params.PageSize = 20
	}

	foundService := new(service.FoundPetService)
	list, total, err := foundService.GetPetList(params)

	if err != nil {
		help.Output(err.Error(), c)
		return
	}
	nowTotal := int64(params.PageSize * params.Page)
	var hasMore bool = true
	if nowTotal >= total {
		hasMore = false
	}
	data := map[string]interface{}{
		"list":      list,
		"total":     total,
		"page_size": params.PageSize,
		"has_more":  hasMore,
		"bannerImage": []string{
			"http://image.foundall.cn/static/gonggao.png",
			"http://image.foundall.cn/static/gonggao2.png",
		},
	}
	help.Output(data, c)
	return

}

// personDetail
//  @Description: 详情页
//  @receiver p
//  @param c
func (p *FoundPetController) Detail(c *gin.Context) {
	var params request.PetDetailParams

	_ = c.ShouldBindBodyWith(&params, binding.JSON)
	foundService := new(service.FoundPetService)
	detail, err := foundService.GetPersonDetail(params.Id)
	if err != nil {
		help.Output(err.Error(), c)
		return
	}
	help.Output(detail, c)
	return
}
