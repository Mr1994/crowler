package request

import "api_client/model"

type PetListParams struct {
	model.CrPet
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

//
//  PersonDetailParams
//  @Description: 信息详情页参数
//
type PertDetailParams struct {
	Id int `json:"id"`
}

//
//  PersonDetailParams
//  @Description: 信息详情页参数
//
type PetDetailParams struct {
	Id int `json:"id"`
}
