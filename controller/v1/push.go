package v1

import (
	"api_client/help"
	middleWare "api_client/middleware"
	"api_client/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"reflect"
)

type te struct {
}

// 存储push数据
func SavePush(c *gin.Context) {
	// 解析到结构体
	var Mess = &service.KafkaData
	json.Unmarshal(MustMarshal(middleWare.GetParams(c)), Mess)
	// 设置Params默认值
	common := new(help.Common)
	common.SetDefault(reflect.TypeOf(Mess).Elem(), reflect.ValueOf(Mess).Elem())

	if Mess.UniqueIdArr == nil {
		panic("请求设备参数为空")
	}
	// 获取跳转的url
	service.KafkaData.SetSchemeUrl()
	// 设置topic
	service.KafkaData.SetTopic()
	// 数据加入kafka
	service.KafkaData.EnqueueWithAlloc()

	help.Output(service.KafkaData, c)
	//new(help.Common).Log("tewrtew", "testfle")
	//f := config.RDS.Get()
	//defer f.Close()
	//f.Do("SET", "hhhh", "value")

}

func MustMarshal(v interface{}) []byte {
	content, _ := json.Marshal(v)
	return content
}
