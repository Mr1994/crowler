package model

// reg_id表
type YwPushRegId struct {
	Id             int64  `json:"id"`              // 自增ID
	UniqueId       string `json:"unique_id"`       // cj_device_info.unique_id,设备的唯一标识
	AppId          int    `json:"app_id"`          // 居理对app的编号
	Brand          string `json:"brand"`           // 手机品牌
	Type           int    `json:"type"`            // 1:小米,2:华为,3:极光
	RegId          string `json:"reg_id"`          // 各厂商生成的regId
	CreateDatetime int    `json:"create_datetime"` // 创建时间
	UpdateDatetime int    `json:"update_datetime"` // 更新时间
	UniqueIdType   int    `json:"unique_id_type"`  // 1:服务端生成,2:客端生成的
	ValidStatus    int    `json:"valid_status"`    // 1 有效 2 无效
}

func (YwPushRegId) TableName() string {
	return "yw_push_reg_id"
}
