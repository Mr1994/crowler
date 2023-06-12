package model

import "api_client/config"

// 期货表
type WtiFutures struct {
	Id           int    `json:"id"`            // 主键id
	Name         string `json:"name"`          // 期货名称
	SerialNumber string `json:"serial_number"` // 期货编号
}

func (w *WtiFutures) TableName() string {
	return "wti_futures"
}

// 查询多条数据
func (w *WtiFutures) Find() []WtiFutures {
	var db = config.DB
	var futrues []WtiFutures
	err := db.Find(&futrues).Error
	if err != nil {
		panic("查询数据异常:" + err.Error())
	}
	return futrues
}
