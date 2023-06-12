package request

type HomeParams struct {
	NativePlace string `json:"native_place"`
}

//
//  HomeSelectParams
//  @Description:
//
type HomeSelectParams struct {
	Page        int    `json:"page"`
	PageSize    int    `json:"page_size"`
	NativePlace string `json:"native_place" gorm:"comment:失踪人口籍贯"`
	MissArea    string `json:"miss_area" gorm:"comment:宠物失踪地"`
}
