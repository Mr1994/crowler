package v1

import (
	"api_client/request"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type LoginController struct {
}

func (p *LoginController) Login(c *gin.Context) {
	// 获取请求参数
	var params request.HomeParams

	_ = c.ShouldBindBodyWith(&params, binding.JSON)

	return

}
