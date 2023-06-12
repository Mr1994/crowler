package model

import (
	"api_client/config"
)

// 期货关联表

type WtiFuturesAssociate struct {
	Id              int    `json:"id"`               // 主键id
	SerialNumber    string `json:"serial_number"`    // 期货编号
	AssociateSerial string `json:"associate_serial"` // 关联id
	Name            string `json:"name"`             // 期货名称
}

func (c *WtiFuturesAssociate) TableName() string {
	return "wti_futures_associate"
}

// 数据查单条数据
func (w *WtiFuturesAssociate) First(serialNumber string) *WtiFuturesAssociate {
	var db = config.DB
	var WtiFuturesAssociate = &WtiFuturesAssociate{}
	err := db.Where("serial_number=?").First(&WtiFuturesAssociate).Error
	if err != nil {
		panic("查询数据异常:" + err.Error())
	}
	return WtiFuturesAssociate
}

// 查询多条数据
func (w *WtiFuturesAssociate) Find(serialNumber string) []WtiFuturesAssociate {
	var db = config.DB
	var WtiFuturesAssociate []WtiFuturesAssociate
	err := db.Where("associate_serial=?", serialNumber).Find(&WtiFuturesAssociate).Error
	if err != nil {
		panic("查询数据异常:" + err.Error())
	}
	return WtiFuturesAssociate
}
