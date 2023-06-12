package apiRoute

import (
	"api_client/controller/v1"
)

// controller分组
type ApiGroup struct {
	FoundPerson v1.FoundPersonController
	Home        v1.HomeController
	FoundPet    v1.FoundPetController
	Upload      v1.UploadController
}

var ApiGroupList = new(ApiGroup)
