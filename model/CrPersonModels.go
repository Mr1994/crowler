package model

import (
	"api_client/config"
	"errors"
	"gorm.io/gorm"
	"time"
)

type CrPerson struct {
	Id             int    `json:"id"`              // 主键id
	Name           string `json:"name"`            // 用户名称
	UserHeader     string `json:"user_header"`     // 用户名称
	MissAddress    string `json:"miss_address"`    // 失踪地址
	Sex            int    `json:"sex"`             // 性别 1 男 2女
	Birthday       int    `json:"birthday"`        // 生日
	MissDay        int    `json:"miss_day"`        // 生日
	NativePlace    string `json:"native_place"`    // 籍贯
	MissHeight     int    `json:"miss_height"`     // 失踪时身高
	Type           int    `json:"type"`            // 寻亲类别： 1家人寻亲
	CallPolice     int    `json:"call_police"`     // 是否报案 1 是 2否
	Relationship   string `json:"relationship"`    // 与被寻人关系 1 家人
	ContactName    string `json:"contact_name"`    // 联系人
	ContactPhone   string `json:"contact_phone"`   // 联系人电话
	MissDetail     string `json:"miss_detail "`    // 失踪详情
	Source         int    `json:"source"`          // 数据来源 0 默认 1. 寻人网 2 希望寻人网
	CrawlerId      string `json:"crawler_id"`      // 抓取编号
	CrawlerAddress string `json:"crawler_address"` // 抓取编号
	PushTime       int    `json:"push_time"`       // 抓取编号
	CreateTime     int64  `json:"create_time"`     // 创建时间
	UpdateTime     int64  `json:"update_time"`     // 更新时间
}

type CrPersonFound struct {
	Id             int    `json:"id"`              // 主键id
	Name           string `json:"name"`            // 用户名称
	UserHeader     string `json:"user_header"`     // 用户名称
	MissAddress    string `json:"miss_address"`    // 失踪地址
	Sex            int    `json:"sex"`             // 性别 1 男 2女
	Birthday       string `json:"birthday"`        // 生日
	MissDay        string `json:"miss_day"`        // 生日
	NativePlace    string `json:"native_place"`    // 籍贯
	MissHeight     int    `json:"miss_height"`     // 失踪时身高
	Type           int    `json:"type"`            // 寻亲类别： 1家人寻亲
	CallPolice     int    `json:"call_police"`     // 是否报案 1 是 2否
	Relationship   string `json:"relationship"`    // 与被寻人关系 1 家人
	ContactName    string `json:"contact_name"`    // 联系人
	ContactPhone   string `json:"contact_phone"`   // 联系人电话
	MissDetail     string `json:"miss_detail"`     // 失踪详情
	Source         int    `json:"source"`          // 数据来源 0 默认 1. 寻人网 2 希望寻人网
	CrawlerId      string `json:"crawler_id"`      // 抓取编号
	CrawlerAddress string `json:"crawler_address"` // 抓取编号
	PushTime       string `json:"push_time"`       // 抓取编号
	CreateTime     string `json:"create_time"`     // 创建时间
	UpdateTime     string `json:"update_time"`     // 更新时间
}

func (c *CrPerson) TableName() string {
	return "cr_person"
}
func (c *CrPerson) Create(CrPersionInfo *CrPerson, All string) bool {
	var db = config.DB
	var crPerson CrPerson

	err := db.Model(CrPerson{}).Where("source=? and crawler_id=?", CrPersionInfo.Source, CrPersionInfo.CrawlerId).First(&crPerson).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 数据为空的处理。
		CrPersionInfo.CreateTime = time.Now().Unix()
		CrPersionInfo.UpdateTime = time.Now().Unix()
		db.Create(&CrPersionInfo) // newUser not user
		return true
	} else {
		// 如果是全量更新的话， 需要更新这条数据
		if All == "true" {
			db.Model(CrPerson{}).Where("source=? and crawler_id=?", CrPersionInfo.Source, CrPersionInfo.CrawlerId).Updates(CrPersionInfo)
		}
		return false
	}
}
