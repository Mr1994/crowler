package middleWare

import (
	"api_client/help"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AllBody struct {
	Common Common
	Params map[string]interface{}
}

// 设置公共请求参数
type Common struct {
	AppId        string `json:"app_id"`
	VersionName  string `json:"version_name"`
	VersionBuild string `json:"version_build"`
	ClientAgent  string `json:"clientagent"`
	DistinctId   string `json:"distinct_id"`
	MjbAppid     int    `json:"mjb_appid"`
	CityId       int    `json:"city_id"`
	Token        string `json:"token"`
	AppVersion   string `json:"version_code"`
}

var postData AllBody

func getPost(request *gin.Context) AllBody {

	err := request.ShouldBindBodyWith(&postData, binding.JSON)
	if err != nil {
		// 返回错误信息
		// gin.H封装了生成json数据的工具
		help.OutError(request, err.Error())
		//request.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		request.Abort()
	}

	return postData
}

// 设置公共参数
func SetCommon(data AllBody, request *gin.Context) {
	// 设置属性
	request.Set("app_id", data.Common.AppId)
	request.Set("clientagent", data.Common.ClientAgent)
	request.Set("city_id", data.Common.CityId)
	request.Set("app_version", data.Common.AppVersion)
	request.Set("token", data.Common.Token)
	request.Set("common", data.Common)
}

// 设置请求参数
func SetParams(data AllBody, request *gin.Context) {
	request.Set("params", data.Params)
}

// 获取客户端的请求参数
func GetParams(request *gin.Context) map[string]interface{} {
	params, err := request.Get("params")
	if err != true {
		panic("参数解析异常")
	}
	post := params.(map[string]interface{})
	return post
}

// 获取map数据
func GetValue(post map[string]interface{}, key string, defaultVallue ...interface{}) interface{} {

	if _, ok := post[key]; ok {
		return post[key]
	} else {
		if len(defaultVallue) == 0 {
			return ""
		} else {
			return defaultVallue[0]
		}
	}
}

// 定义中间 不生效
func RequestMiddleWare() gin.HandlerFunc {
	return func(request *gin.Context) {
		// 获取请求参数
		//data :=
		//(request)
		//content, _ := json.Marshal(data)
		//fmt.Println(string(content))
		// 设置请求属性
		// TODO 将数据存储到key中 但是如果这么设置，如果获取common这个key下面的其他属性
		// 获取请求基本参数
		//SetCommon(data, request)
		// 设置请求参数
		//SetParams(data, request)
		defer func() {
			if err := recover(); err != nil {
				help.OutError(request, err.(string))
				return
			}

		}()
		// 如果有错误则直接获取，并return' 一个固定结构
		request.Next()
	}
}
