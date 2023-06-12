package v1

import (
	"api_client/constant"
	"api_client/help"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"path"
	"strings"
)

type UploadController struct {
}

/*
上传图片路径
*/
const uploadPath = "source/image/own/"

// personDetail
//  @Description: 上传图片
//  @receiver p
//  @param c
func (p *UploadController) UploadFile(c *gin.Context) {
	// 接收图片
	f, err := c.FormFile("imgfile")
	if err != nil {
		help.OutError(c, err.Error())
		return
	} else {

		fileExt := strings.ToLower(path.Ext(f.Filename))
		if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" {
			help.OutError(c, "上传失败!只允许png,jpg,gif,jpeg文件")
			return
		}
		tools := new(help.Tools)
		fileName := tools.GenerateMd5(f.Filename)
		dirPath := "./" + uploadPath
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			// mkdir 创建目录，mkdirAll 可创建多层级目录
			os.MkdirAll(dirPath, os.ModePerm)
		}
		fildDir := fmt.Sprintf(dirPath)
		filepath := fmt.Sprintf("%s%s%s", fildDir, fileName, fileExt)
		c.SaveUploadedFile(f, filepath)
		var pathUrl map[string]string
		pathUrl = make(map[string]string)
		pathUrl["path"] = constant.ImageUrl + filepath
		help.Output(pathUrl, c)

	}
}

//判断文件文件夹是否存在
func isFileExist(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, nil
	}
	//我这里判断了如果是0也算不存在
	if fileInfo.Size() == 0 {
		return false, nil
	}
	if err == nil {
		return true, nil
	}
	return false, err
}
