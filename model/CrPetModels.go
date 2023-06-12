package model

import (
	"api_client/config"
	"fmt"
	"time"
)

type CrPet struct {
	Id             int    `json:"id"`              // 主键
	Title          string `json:"title"`           // 文章标题
	MissDay        int    `json:"miss_day"`        // 走失日期
	Source         int    `json:"source"`          // 1 寻狗网 （www.xungou5.com）
	Breed          string `json:"breed"`           // 品种
	ContractPhone  string `json:"contract_phone"`  // 联系电话
	MissDetail     string `json:"miss_detail"`     // 丢失详情
	CrawlerId      string `json:"crawler_id"`      // 抓取编号
	CrawlerAddress string `json:"crawler_address"` // 抓取编号
	MissArea       string `json:"miss_area"`       // 丢失区域
	FoundType      int    `json:"found_type"`      // 专区类型，1 寻找狗狗
	CreateTime     int64  `json:"create_time"`     // 创建时间
	UpdateTime     int64  `json:"update_time"`     // 更新时间
}

func (c *CrPet) TableName() string {
	return "cr_pet"
}
func (c *CrPet) Create(CrPetInfo *CrPet, All string) bool {
	var db = config.DB
	var crPet CrPet

	err := db.Model(CrPet{}).Where("source=? and crawler_id=?", CrPetInfo.Source, CrPetInfo.CrawlerId).First(&crPet).Error

	// 判断是否有重复数据
	if err != nil && err.Error() == "record not found" {
		fmt.Println(111111)
		// 数据为空的处理。
		CrPetInfo.CreateTime = time.Now().Unix()
		CrPetInfo.UpdateTime = time.Now().Unix()
		db.Create(&CrPetInfo) // newUser not user
		return true
	} else {
		fmt.Println(2222)
		// 如果是全量更新的话， 需要更新这条数据
		if All == "true" {
			db.Model(CrPet{}).Where("source=? and crawler_id=?", CrPetInfo.Source, CrPetInfo.CrawlerId).Updates(CrPetInfo)
		}
		return false
	}
}
