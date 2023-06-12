package service

import (
	"api_client/config"
	"api_client/constant"
	"api_client/help"
	"api_client/model"
	v1 "api_client/request"
	"encoding/json"
	"fmt"
	"strings"
)

/**
爬虫页面结构
*/
type FoundPetService struct {
}

type CrPetList struct {
	Id             int                    `json:"id"`             // 主键
	Title          string                 `json:"title"`          // 文章标题
	MissDay        string                 `json:"miss_day"`       // 走失日期
	Breed          string                 `json:"breed"`          // 品种
	ContractPhone  string                 `json:"contract_phone"` // 联系电话
	MissArea       string                 `json:"miss_area"`      // 丢失区域
	FoundType      int                    `json:"found_type"`     // 专区类型，1 寻找狗狗
	CreateTime     string                 `json:"create_time"`    // 创建时间
	MissDetail     string                 `json:"miss_detail"`    // 丢失详情
	UpdateTime     string                 `json:"update_time"`    // 更新时间
	MissDetailInfo map[string]interface{} `json:"miss_detail_info"`
	UserHeader     string                 `json:"user_header"`
}

func AssemblyCrawlerPet(text string, crPetModel *model.CrPet) {
	/**
	  寻人编号
	*/
	var matchingString []string
	matchingString = make([]string, 4)
	matchingString = []string{"区域", "年份", "品种", "类型", "日期", "联系"}
	var common *help.Common // 定义类型
	common = new(help.Common)
	Cps := new(FoundPetService)
	for _, v := range matchingString {
		hasString := strings.Index(text, v)
		// 如果匹配到了数据，则看是哪一个数据
		if hasString >= 0 {
			runeString := []byte(text)
			getStringPos := strings.Index(text, "：")
			matchingElementByte := runeString[getStringPos+3 : len(text)]
			// 转化成字符查找空格
			matchingElement := string(matchingElementByte)
			// 查找第一个空格出现时间
			getEmptyStringPos := strings.Index(matchingElement, " ")
			// 从第一个空格开始之后， 就不要后面的数据了
			if getEmptyStringPos >= 0 {
				matchingElement = string(matchingElementByte[0:getEmptyStringPos])
			}
			switch v {
			case "区域":
				crPetModel.MissArea = matchingElement
			case "日期":
				crPetModel.MissDay = common.GetTimeStamp(matchingElement, "2006-01-02")
			case "品种":
				crPetModel.Breed = matchingElement
			case "类型":
				crPetModel.FoundType = Cps.GetFoundType(matchingElement)
			case "联系":
				// 如果如果前面没有抓取到手机号码，这里补充抓取
				if matchingElement != "" && crPetModel.ContractPhone == "" {
					isMobile := common.CheckMobile(matchingElement)
					fmt.Println(isMobile)
					if isMobile {
						crPetModel.ContractPhone = matchingElement
					}
				}

			}
		}
	}
}

func (c *FoundPetService) GetFoundType(foundString string) int {

	findType, ok := constant.FoundType[foundString]
	if ok {
		return findType
	} else {
		return 0 // 没匹配默认其他
	}

}

func (c *FoundPetService) GetHomePetList(info *v1.HomeSelectParams) (list interface{}, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var dbConnect = config.DB.Debug()

	// 创建db
	db := dbConnect.Model(model.CrPet{})
	var PetList []*CrPetList
	// 如果有条件搜索 下方会自动创建搜索语句

	if info.MissArea != "" {
		db = db.Where("native_place LIKE ?", "%"+info.MissArea+"%")
	}

	err = db.Limit(limit).Offset(offset).Order("id desc").Find(&PetList).Error

	if err != nil {
		return
	}
	var common *help.Common
	common = new(help.Common)
	for _, v := range PetList {
		v.MissDay = common.GetFormatTime(v.MissDay)
		v.CreateTime = common.GetFormatTime(v.CreateTime)

	}
	return PetList, err
}

// 获取失踪人口列表
func (c *FoundPetService) GetPetList(info *v1.HomeSelectParams) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var dbConnect = config.DB.Debug()

	// 创建db
	db := dbConnect.Model(model.CrPet{})
	var PetList []*CrPetList
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.MissArea != "" {
		db = db.Where("name LIKE ?", "%"+info.MissArea+"%")
	}
	if info.NativePlace != "" {
		db = db.Where("native_place LIKE ?", "%"+info.NativePlace+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Order("id desc").Find(&PetList).Error
	common := new(help.Common)
	// 格式化时间
	for _, v := range PetList {
		v.MissDay = common.GetFormatTime(v.MissDay)
		v.CreateTime = common.GetFormatTime(v.CreateTime)

		var result map[string]interface{}
		if err := json.Unmarshal([]byte(v.MissDetail), &result); err != nil {
			v.UserHeader = ""
			continue
		}
		v.MissDetailInfo = result
		if _, ok := v.MissDetailInfo["img"]; ok {
			UserHeaderList := v.MissDetailInfo["img"].([]interface{})

			if UserHeaderList != nil {
				v.UserHeader = UserHeaderList[0].(string)
			}
		} else {
			v.UserHeader = ""
		}
	}
	return PetList, total, err
}

// 获取失踪人口列表
func (c *FoundPetService) GetPersonDetail(Id int) (list interface{}, err error) {
	var dbConnect = config.DB.Debug()

	// 创建db
	db := dbConnect.Model(model.CrPet{})
	var PetList *CrPetList
	// 如果有条件搜索 下方会自动创建搜索语句

	err = db.Where("id = ?", Id).Order("id desc").Find(&PetList).Error
	common := new(help.Common)
	// 格式化时间
	PetList.MissDay = common.GetFormatTime(PetList.MissDay)
	PetList.CreateTime = common.GetFormatTime(PetList.CreateTime)
	if PetList.MissDetail != "" {
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(PetList.MissDetail), &result); err != nil {
			return nil, err
		}
		PetList.MissDetailInfo = result
	}
	return PetList, err
}
