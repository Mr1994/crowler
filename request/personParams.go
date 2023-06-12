package request

import "api_client/model"

type PersonListParams struct {
	model.CrPerson
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

//
//  PersonDetailParams
//  @Description: 信息详情页参数
//
type PersonDetailParams struct {
	Id int `json:"id"`
}
