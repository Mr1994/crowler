package service

import (
	"api_client/config"
	"api_client/constant"
	"api_client/help"
	"api_client/model"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

const hasFilter = 2 // 是否需要对数据进行过滤
const priority = 1  // push优先级
const sendType = 0  // push类型

type KafkaMessage struct {
	JuliveAppId  int        `json:"julive_app_id"`
	HasFilter    int        `json:"has_filter"  default:"1"`
	RequestId    string     `json:"request_id"`
	CreateTime   int64      `json:"create_time"`
	UniqueIdArr  []string   `json:"unique_id_arr"`
	Title        string     `json:"title"`
	Notification string     `json:"notification"`
	Batch        string     `json:"batch"`
	PushConfig   PushConfig `json:"push_config"`
	PushParams   PushParams `json:"push_params"`
	PushTime     string     `json:"timestamp"`
	Priority     int        `json:"priority" default:"1"`
	SendType     int        `json:"send_type" default:"0"`
	regIdArr     map[string]string
}
type PushParams struct {
	SchemeUrl string `json:"scheme_url"`
}

type PushConfig struct {
	UserId          int `json:"user_id"`
	OperationUserId int `json:"operation_user_id"`
	CommendId       int `json:"commend_id"`
	InformationId   int `json:"information_id"`
	Level           int `json:"level"`
	Title           string
	Notification    string
	Behavior        int
	PushType        int `json:"push_type"`
	Devicetype      int `json:"type"`
	TopicType       int `json:"topic_type"`
	PassThrough     int `json:"passThrough"`
	PushNow         int `json:"push_now"`
}

var KafkaData KafkaMessage

var appH5SchemeUrl = map[int]string{
	constant.AppIos:     "comjia://app.comjia.com/h5?data=%s",
	constant.AppAndroid: "comjia://app.comjia.com/h5?data=%s",
}

// 设置发送类型
var kafkaSendTypeTopic = map[int]string{
	constant.SendTypeAuto:       "Julive_Queue_Push_Service_JIGUANG",
	constant.SendTypeXiaomi:     "Julive_Queue_Push_Service_XIAOMI",
	constant.SendTypeHuawei:     "Julive_Queue_Push_Service_HUAWEI",
	constant.SendTypeJiguang:    "Julive_Queue_Push_Service_JIGUANG",
	constant.SendTypeVivo:       "Julive_Queue_Push_Service_VIVO",
	constant.SendTypeOppo:       "Julive_Queue_Push_Service_OPPO",
	constant.SendTypeJiguangVip: "Julive_Queue_Push_Service_JIGUANG",
	constant.SendTypeApple:      "Julive_Queue_Push_Service_APNS",
}

// 设置跳转协议
func (k *KafkaMessage) SetSchemeUrl() {
	if KafkaData.PushParams.SchemeUrl == "" {
		panic("跳转协议为空")
	}
	// 如果是http协议，直接返回协议地址
	if strings.Index(KafkaData.PushParams.SchemeUrl, "http") == 0 {
		return
	}

	// 对url进行编码
	var schemeParams = make(map[string]interface{})
	schemeParams["url"] = KafkaData.PushParams.SchemeUrl
	String, _ := json.Marshal(schemeParams)
	jsonString := url.QueryEscape(string(String))
	schemeUrl := fmt.Sprintf(appH5SchemeUrl[KafkaData.JuliveAppId], jsonString)
	KafkaData.PushParams.SchemeUrl = schemeUrl
}

// 设置push的topic
func (k *KafkaMessage) SetTopic() string {
	if _, ok := kafkaSendTypeTopic[KafkaData.SendType]; ok {
		return kafkaSendTypeTopic[KafkaData.SendType]
	}
	return kafkaSendTypeTopic[constant.SendTypeAuto]

}

// 设置push的topic
func (k *KafkaMessage) SetGroup() string {
	if _, ok := constant.PriorityMapGroupId[KafkaData.Priority]; ok {
		return constant.PriorityMapGroupId[KafkaData.Priority]
	}
	return constant.PriorityMapGroupId[constant.PriorityLevel1]

}

func (k *KafkaMessage) EnqueueWithAlloc() bool {
	var regIdList map[int]map[string]string
	regIdList = k.getRegIdList()
	if regIdList == nil {
		new(help.Common).Log("未找到设备信息", "pushInfo")
		return true
	}
	for k, v := range regIdList {
		KafkaData.SendType = k
		KafkaData.regIdArr = v
		KafkaData.UniqueIdArr = nil
		for k1, _ := range v {
			KafkaData.UniqueIdArr = append(KafkaData.UniqueIdArr, k1)
		}
	}
	//new(help.Common).PushSyncKafkaMessage("k2_test3", KafkaData)

	new(help.Common).PushAsyncKafkaMessage("k2_test3", KafkaData)
	return true

}

// 查询推送过来的设备的id
func (k *KafkaMessage) getRegIdList() map[int]map[string]string {
	db := config.DB

	var ywPushRegId []model.YwPushRegId
	var err = db.Debug().Where("app_id = ? and unique_id in (?) and valid_status = ? ", KafkaData.JuliveAppId, KafkaData.UniqueIdArr, 1).Find(&ywPushRegId).Error
	if err != nil {
		fmt.Println(err)
	}

	//var regIdList map[int]map[string]string
	regIdList := make(map[int]map[string]string)
	uniqueMap := make(map[string]string)
	for _, v := range ywPushRegId {
		uniqueId := v.UniqueId
		typeId := v.Type
		RegId := v.RegId
		if _, ok := uniqueMap[uniqueId]; !ok {
			uniqueMap[uniqueId] = RegId
			if regIdList[typeId] == nil {
				regIdList[typeId] = make(map[string]string)
			}
			regIdList[typeId][uniqueId] = RegId
		}
	}

	return regIdList

}
