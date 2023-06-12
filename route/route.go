package route

import (
	"api_client/controller/apiRoute"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine) {

	// v1 版本
	GroupV1 := r.Group("/v1")
	var apiGroup = apiRoute.ApiGroupList
	{
		// foundPersonController
		GroupV1.POST("foundPerson/personList", apiGroup.FoundPerson.List)
		GroupV1.POST("foundPerson/personDetail", apiGroup.FoundPerson.PersonDetail)
		GroupV1.POST("foundPerson/submit", apiGroup.FoundPerson.AddPerson)

		// HomeController
		GroupV1.POST("home/homeInfo", apiGroup.Home.HomeInfo)

		// PetController
		GroupV1.POST("foundPet/petList", apiGroup.FoundPet.List)
		GroupV1.POST("foundPet/petDetail", apiGroup.FoundPet.Detail)
		// uploadController
		GroupV1.POST("upload/uploadFile", apiGroup.Upload.UploadFile)

	}

}
