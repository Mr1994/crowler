package model

import (
	"api_client/config"
	"bytes"
	"fmt"
	"gorm.io/gorm"
)

type WtiFuturesStatic struct {
	Id             int    `json:"id"`               // 主键id
	Name           string `json:"name"`             // 期货名称
	SerialNumber   string `json:"serial_number"`    // 期货编号
	Date           string `json:"date"`             // 日期
	MaxPrice       string `json:"max_price"`        // 当日最高价格
	MinPrice       string `json:"min_price"`        // 当日最低价格
	OpenPrice      string `json:"open_price"`       // 当日最低价格
	SpotGoodsName  string `json:"spot_goods_name"`  // 现货名称
	SpotGoodsPrice string `json:"spot_goods_price"` // 现货价格
}

func (c *WtiFuturesStatic) TableName() string {
	return "wti_futures_static"
}

// BatchSave 批量插入数据
func (c *WtiFuturesStatic) BatchSave(futuresList []WtiFuturesStatic) error {
	var db = config.DB
	var buffer bytes.Buffer
	sql := "insert into " + (&WtiFuturesStatic{}).TableName() + " (`name`,`serial_number`,`date`,`max_price`,`min_price`,`open_price`,`spot_goods_name`,`spot_goods_price`) values"
	if _, err := buffer.WriteString(sql); err != nil {
		return err
	}
	for i, w := range futuresList {
		if i == len(futuresList)-1 {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s','%s','%s','%s','%s','%s');", w.Name, w.SerialNumber, w.Date, w.MaxPrice, w.MinPrice, w.OpenPrice, w.SpotGoodsName, w.SpotGoodsPrice))
		} else {
			buffer.WriteString(fmt.Sprintf("('%s','%s','%s','%s','%s','%s','%s','%s'),", w.Name, w.SerialNumber, w.Date, w.MaxPrice, w.MinPrice, w.OpenPrice, w.SpotGoodsName, w.SpotGoodsPrice))
		}
	}
	return db.Exec(buffer.String()).Error
}

// 更新数据存储
func (c *WtiFuturesStatic) Create(FuturesStatic *WtiFuturesStatic, ch chan int) {
	var futureStaci = WtiFuturesStatic{}
	var db = config.DB
	err := db.Where("date=? and serial_number=?", FuturesStatic.Date, FuturesStatic.SerialNumber).First(&futureStaci).Error

	// error handling...
	if gorm.ErrRecordNotFound == err {
		db.Create(&FuturesStatic) // newUser not user
	} else {
		db.Model(&FuturesStatic).Where("date=? and serial_number=?", FuturesStatic.Date, FuturesStatic.SerialNumber).Updates(&FuturesStatic)
	}
	<-ch
}

// 查询多条数据
func (w *WtiFuturesStatic) Find(serialNumber string) []WtiFuturesStatic {
	var db = config.DB
	var futrues []WtiFuturesStatic
	err := db.Find(&futrues).Where("serial_number=?", serialNumber).Error
	if err != nil {
		panic("查询数据异常:" + err.Error())
	}
	return futrues
}
