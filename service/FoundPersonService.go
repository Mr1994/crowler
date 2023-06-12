package service

import (
	"api_client/config"
	"api_client/constant"
	"api_client/help"
	"api_client/model"
	v1 "api_client/request"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

/**
爬虫页面结构
*/
type FoundPersonService struct {
}

type CrSelectPerson struct {
	Id             int                    `json:"id"`              // 主键id
	Name           string                 `json:"name"`            // 用户名称
	UserHeader     string                 `json:"user_header"`     // 用户名称
	MissAddress    string                 `json:"miss_address"`    // 失踪地址
	Sex            int                    `json:"sex"`             // 性别 1 男 2女
	Birthday       string                 `json:"birthday"`        // 生日
	MissDay        string                 `json:"miss_day"`        // 生日
	NativePlace    string                 `json:"native_place"`    // 籍贯
	MissHeight     int                    `json:"miss_height"`     // 失踪时身高
	Type           int                    `json:"type"`            // 寻亲类别： 1家人寻亲
	CallPolice     int                    `json:"call_police"`     // 是否报案 1 是 2否
	Relationship   string                 `json:"relationship"`    // 与被寻人关系 1 家人
	ContactName    string                 `json:"contact_name"`    // 联系人
	ContactPhone   string                 `json:"contact_phone"`   // 联系人电话
	MissDetail     string                 `json:"miss_detail "`    // 失踪详情
	Source         int                    `json:"source"`          // 数据来源 0 默认 1. 寻人网 2 希望寻人网
	CrawlerId      string                 `json:"crawler_id"`      // 抓取编号
	CrawlerAddress string                 `json:"crawler_address"` // 抓取编号
	PushTime       string                 `json:"push_time"`       // 抓取编号
	CreateTime     string                 `json:"create_time"`     // 创建时间
	UpdateTime     string                 `json:"update_time"`     // 更新时间
	MissDetailInfo map[string]interface{} `json:"miss_detail_info"`
	UserHeaderList []string               `json:"user_header_list"` // 用户名称

}

type UserHeaderUpdate struct {
	UserHeader string `json:"user_header"` // 用户名称
}

/**
组装爬虫数据(寻人网) http://www.xunrenla.com/
*/
func (c *FoundPersonService) AssemblyCrawlerData(crawlerBody string, cpPersonStruct *model.CrPerson) {

	// 失踪地址
	c.getCrawlerMissAddress(crawlerBody, cpPersonStruct)
	// 身高
	c.getCrawlerMissHeight(crawlerBody, cpPersonStruct)
	// 类型
	c.getCrawlerType(crawlerBody, cpPersonStruct)
	// 报案
	c.getCrawlerCallPolice(crawlerBody, cpPersonStruct)
	// 联系人
	c.getCrawlerContactName(crawlerBody, cpPersonStruct)
	// 联系电话
	c.getCrawlerContactPhone(crawlerBody, cpPersonStruct)
	// 与被寻人的关系
	c.getCrawlerRelationship(crawlerBody, cpPersonStruct)

}

/*
获取用户性别
*/
func (c *FoundPersonService) getCrawlerSex(crawlerBody string, cpPersonStruct *model.CrPerson) {
	strIndex := strings.Index(crawlerBody, "性别")
	// 判断是否存在姓名
	if strIndex >= 0 {
		sex := strings.Replace(crawlerBody, "性别：", "", 1)

		var sexInt int
		if sex == "男" {
			sexInt = constant.Man
		} else {
			sexInt = constant.WoMan

		}
		cpPersonStruct.Sex = sexInt
	}
}

/**
获取用户名称
*/
func (c *FoundPersonService) getCrawlerName(crawlerBody string, cpPersonStruct *model.CrPerson) {
	strIndex := strings.Index(crawlerBody, "姓名")
	// 如果存在姓名字符，则添加到结构体中
	if strIndex >= 0 {
		name := strings.Replace(crawlerBody, "姓名：", "", 1)
		cpPersonStruct.Name = name

	}
}

/**
生日
*/
func (c *FoundPersonService) getCrawlerBirthDay(crawlerBody string, cpPersonStruct *model.CrPerson) {
	strIndex := strings.Index(crawlerBody, "出生日期")
	common := new(help.Common)
	// 时间格式如果异常，则不显示
	defer func() {
		err := recover()
		if err != nil {
			cpPersonStruct.Birthday = 0
			new(help.Common).Log(err.(string), "cmd")

		}
	}()

	if strIndex >= 0 {
		matchStr := strings.Replace(crawlerBody, "出生日期：", "", 1)
		var birthDay int
		// 不同网站时间格式类型不一致
		birthDay = common.GetTimeStamp(matchStr, "2006-01-02")
		cpPersonStruct.Birthday = birthDay
	}
}
func (c *FoundPersonService) GetCrawlerPushTime(crawlerBody string, cpPersonStruct *model.CrPerson) {
	strIndex := strings.Index(crawlerBody, "发布时间")
	common := new(help.Common)
	// 时间格式如果异常，则不显示
	defer func() {
		err := recover()
		if err != nil {
			cpPersonStruct.PushTime = 0
			new(help.Common).Log(err.(string), "cmd")

		}
	}()

	if strIndex >= 0 {
		matchStr := strings.Replace(crawlerBody, "发布时间：", "", 1)
		var pushTime int
		// 不同网站时间格式类型不一致
		pushTime = common.GetTimeStamp(matchStr, "2006-01-02 15:04")
		cpPersonStruct.PushTime = pushTime
	}
}

/**
失踪日期
*/
func (c *FoundPersonService) getCrawlerMissDay(crawlerBody string, cpPersonStruct *model.CrPerson) {

	// 不同网站匹配的名字不一样
	var strRep string
	if cpPersonStruct.Source == constant.SourceHopeXunRen {
		strRep = "失散日期"
	} else {
		strRep = "失踪日期"
	}
	strIndex := strings.Index(crawlerBody, strRep)
	common := new(help.Common)
	// 没有获取到时间，则不显示即可
	defer func() {
		err := recover()
		if err != nil {
			new(help.Common).Log(err.(string), "cmd")
			cpPersonStruct.MissDay = 0
		}
	}()
	// 如果存在姓名字符，则添加到结构体中
	if strIndex >= 0 {
		matchStr := strings.Replace(crawlerBody, strRep+"：", "", 1)
		var missDay int
		missDay = common.GetTimeStamp(matchStr, "2006-01-02")
		cpPersonStruct.MissDay = missDay
	}

}

/**
籍贯
*/
func (c *FoundPersonService) getCrawlerNativePlace(crawlerBody string, cpPersonStruct *model.CrPerson) {
	strIndex := strings.Index(crawlerBody, "籍贯")
	// 如果存在姓名字符，则添加到结构体中
	if strIndex >= 0 {
		matchStr := strings.Replace(crawlerBody, "籍贯：", "", 1)
		cpPersonStruct.NativePlace = matchStr
	}
}

/**
寻人编号
*/
func (c *FoundPersonService) GetCrawlerId(url string, cpPersonStruct *model.CrPerson) {
	parts := strings.Split(url, "/")
	lastPart := parts[len(parts)-1]                 // 取最后一部分的值
	result := strings.TrimSuffix(lastPart, ".html") // 去掉 .html 后缀
	cpPersonStruct.CrawlerId = result

}

/**
失踪地址
*/
func (c *FoundPersonService) getCrawlerMissAddress(crawlerBody string, cpPersonStruct *model.CrPerson) {
	var strRep string
	if cpPersonStruct.Source == constant.SourceHopeXunRen {
		strRep = "失散地点"
	} else {
		strRep = "失踪地点"
	}
	strIndex := strings.Index(crawlerBody, strRep)

	// 如果存在姓名字符，则添加到结构体中
	if strIndex >= 0 {
		matchStr := strings.Replace(crawlerBody, strRep+"：", "", 1)
		cpPersonStruct.MissAddress = matchStr
	}
}

/**
与被寻人关系
*/
func (c *FoundPersonService) getCrawlerRelationship(crawlerBody string, cpPersonStruct *model.CrPerson) {
	strIndex := strings.Index(crawlerBody, "与被寻者的关系")
	// 如果存在姓名字符，则添加到结构体中
	if strIndex >= 0 {
		matchStr := strings.Replace(crawlerBody, "与被寻者的关系：", "", 1)
		cpPersonStruct.Relationship = matchStr
	}
}

/**
身高
*/
func (c *FoundPersonService) getCrawlerMissHeight(crawlerBody string, cpPersonStruct *model.CrPerson) {
	var strRep string
	if cpPersonStruct.Source == constant.SourceHopeXunRen {
		strRep = "身高"
	} else {
		strRep = "失踪时身高"
	}
	strIndex := strings.Index(crawlerBody, strRep)
	// 如果存在姓名字符，则添加到结构体中
	if strIndex >= 0 {
		matchStr := strings.Replace(crawlerBody, strRep+"：", "", 1)
		matchStr = strings.Replace(matchStr, "CM", "", 1)
		height, _ := strconv.Atoi(matchStr) // 转化成int
		cpPersonStruct.MissHeight = height
	}
}

func (c *FoundPersonService) getCrawlerType(crawlerBody string, cpPersonStruct *model.CrPerson) {
	strIndex := strings.Index(crawlerBody, "寻人类别")
	// 如果存在姓名字符，则添加到结构体中
	if strIndex >= 0 {
		matchStr := strings.Replace(crawlerBody, "寻人类别：", "", 1)
		var missType int
		missType = constant.FindType[strings.TrimSpace(matchStr)]
		cpPersonStruct.Type = missType
	}
}
func (c *FoundPersonService) getCrawlerCallPolice(crawlerBody string, cpPersonStruct *model.CrPerson) {
	strIndex := strings.Index(crawlerBody, "是否报案")
	// 如果存在姓名字符，则添加到结构体中
	if strIndex >= 0 {
		matchStr := strings.Replace(crawlerBody, "是否报案：", "", 1)
		var callPolice int
		var repMatchStr string
		if cpPersonStruct.Source == constant.SourceHopeXunRen {
			repMatchStr = "否"
		} else {
			repMatchStr = "未报案"
		}
		if matchStr == repMatchStr {
			callPolice = constant.CallPoliceNo
		} else {
			callPolice = constant.CallPoliceYes
		}
		cpPersonStruct.CallPolice = callPolice
	}
}

func (c *FoundPersonService) getCrawlerContactName(crawlerBody string, cpPersonStruct *model.CrPerson) {
	strIndex := strings.Index(crawlerBody, "联系人")
	// 如果存在姓名字符，则添加到结构体中
	if strIndex >= 0 {
		matchStr := strings.Replace(crawlerBody, "联系人：", "", 1)
		cpPersonStruct.ContactName = matchStr
	}
}
func (c *FoundPersonService) getCrawlerContactPhone(crawlerBody string, cpPersonStruct *model.CrPerson) {
	strIndex := strings.Index(crawlerBody, "联系电话")
	// 如果存在姓名字符，则添加到结构体中
	if strIndex >= 0 {
		strPhIndex := strings.LastIndex(crawlerBody, ":")
		contractPhone := crawlerBody[strPhIndex+1 : len(crawlerBody)-3]
		cpPersonStruct.ContactPhone = contractPhone

	}
}

// 获取失踪人口列表
func (c *FoundPersonService) GetPersonList(info *v1.PersonListParams) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var dbConnect = config.DB.Debug()

	// 创建db
	db := dbConnect.Model(model.CrPerson{})
	var PersonList []*model.CrPersonFound
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	if info.NativePlace != "" {
		db = db.Where("native_place LIKE ?", "%"+info.NativePlace+"%")
	}
	if info.UserHeader != "" {
		db = db.Where("user_header LIKE ?", "%"+info.UserHeader+"%")

	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("id desc").Find(&PersonList).Error
	common := new(help.Common)
	// 格式化时间
	for _, v := range PersonList {
		v.PushTime = common.GetFormatTime(v.PushTime)
		v.Birthday = common.GetFormatTime(v.Birthday)
		v.MissDay = common.GetFormatTime(v.MissDay)
		v.CreateTime = common.GetFormatTime(v.CreateTime)
		v.UpdateTime = common.GetFormatTime(v.UpdateTime)
		v.UserHeader = common.GetImageUrl(v.UserHeader)

	}
	return PersonList, total, err
}

func (c *FoundPersonService) GetHomePersonList(info *v1.HomeSelectParams) (list interface{}, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var dbConnect = config.DB.Debug()

	// 创建db
	db := dbConnect.Model(model.CrPerson{})
	var PersonList []*CrSelectPerson
	// 如果有条件搜索 下方会自动创建搜索语句

	if info.NativePlace != "" {
		db = db.Where("native_place LIKE ?", "%"+info.NativePlace+"%")
	}

	err = db.Limit(limit).Offset(offset).Order("id desc").Find(&PersonList).Error

	if err != nil {
		return
	}
	fmt.Println(PersonList)
	common := help.Common{}
	for _, v := range PersonList {
		v.PushTime = common.GetFormatTime(v.PushTime)
		v.Birthday = common.GetFormatTime(v.Birthday)
		v.MissDay = common.GetFormatTime(v.MissDay)
		v.CreateTime = common.GetFormatTime(v.CreateTime)
		v.UserHeader = common.GetImageUrl(v.UserHeader)

	}
	return PersonList, err
}

// UpdateDateById
//  @Description: 根据id更新数据
//  @receiver c
//  @param CrPersionInfo
//  @param All
func (c *FoundPersonService) UpdateDateById(Id int, HeaderInfo UserHeaderUpdate) {
	var db = config.DB
	db.Model(model.CrPerson{}).Where("id=?", Id).Updates(HeaderInfo)

}

// GetPersonDetail
//  @Description:  获取用户丢失详情
//  @receiver c
//  @param id
func (c *FoundPersonService) GetPersonDetail(id int) (detail *CrSelectPerson, err error) {
	var dbConnect = config.DB.Debug()

	// 创建db
	db := dbConnect.Model(model.CrPerson{})
	var PersonInfo *CrSelectPerson
	// 如果有条件搜索 下方会自动创建搜索语句})
	err = db.Where("id =?", id).Find(&PersonInfo).Error
	if err != nil {
		return nil, err
	}
	if PersonInfo.MissDetail != "" {
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(PersonInfo.MissDetail), &result); err != nil {
			return nil, err
		}
		PersonInfo.MissDetailInfo = result
	}
	common := help.Common{}
	var userHeaderList []string
	userHeaderArray := make([]string, 3)
	userHeaderArray = append(userHeaderList, common.GetImageUrl(PersonInfo.UserHeader))
	PersonInfo.PushTime = common.GetFormatTime(PersonInfo.PushTime)
	PersonInfo.Birthday = common.GetFormatTime(PersonInfo.Birthday)
	PersonInfo.MissDay = common.GetFormatTime(PersonInfo.MissDay)
	PersonInfo.CreateTime = common.GetFormatTime(PersonInfo.CreateTime)
	PersonInfo.UserHeaderList = userHeaderArray
	return PersonInfo, nil
}

func (c *FoundPersonService) Create(CrPersionInfo *model.CrPerson, All string) bool {
	var db = config.DB
	// 数据为空的处理。
	CrPersionInfo.CreateTime = time.Now().Unix()
	CrPersionInfo.UpdateTime = time.Now().Unix()
	db.Create(&CrPersionInfo) // newUser not user
	return true

}
