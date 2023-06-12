package v1

import (
	"api_client/help"
	"api_client/model"
	"api_client/request"
	"api_client/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type FoundPersonController struct {
}

/**
获取失踪人口用户列表
*/
func (p *FoundPersonController) List(c *gin.Context) {
	// 获取请求参数
	var params request.PersonListParams

	_ = c.ShouldBindBodyWith(&params, binding.JSON)
	if params.PageSize == 0 {
		params.PageSize = 20
	}
	foundService := new(service.FoundPersonService)
	list, total, err := foundService.GetPersonList(&params)

	if err != nil {
		help.Output(err.Error(), c)
		return
	}

	data := map[string]interface{}{
		"list":      list,
		"total":     total,
		"page_size": params.PageSize,
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
func (p *FoundPersonController) PersonDetail(c *gin.Context) {
	var params request.PersonDetailParams

	_ = c.ShouldBindBodyWith(&params, binding.JSON)
	foundService := new(service.FoundPersonService)
	detail, err := foundService.GetPersonDetail(params.Id)
	if err != nil {
		help.Output(err.Error(), c)
		return
	}
	help.Output(detail, c)
	return
}
func (p *FoundPersonController) AddPerson(c *gin.Context) {
	var params *model.CrPerson
	_ = c.ShouldBindBodyWith(&params, binding.JSON)
	CrPersonService := new(service.FoundPersonService)
	CrPersonService.Create(params, "false")

}
