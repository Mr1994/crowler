package cmd

import (
	"api_client/help"
	"api_client/model"
	"api_client/request"
	"api_client/service"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"sync"
)

type FixCommand struct {
}

var FixImageCommand = &cobra.Command{

	Use:   "fixImage",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		dirName, _ := cmd.Flags().GetString("dir")
		web, _ := cmd.Flags().GetString("web")
		helpCommon := new(help.Common)
		pwd, _ := os.Getwd()

		dir := helpCommon.GetParentDirectory(pwd)
		// 查询所有图片数据
		var page int = 0
		httpHelp := new(help.HttpHelp)
		fix := new(FixCommand)
		pwdDir := fix.getDirPath(dir, dirName)

		params := new(request.PersonListParams)
		PersonService := new(service.FoundPersonService)

		foundPersonService := new(service.FoundPersonService)

		count := 10            // 最大支持并发
		wg := sync.WaitGroup{} //控制主协程等待所有子协程执行完之后再退出。

		c := make(chan int, count) // 控制任务并发的chan
		defer close(c)

		for {
			params.UserHeader = web
			params.Page = 0
			params.PageSize = 10
			list, _, _ := PersonService.GetPersonList(params)
			if list == nil {
				panic("数据抓取完成")
			}
			var userHeaderData = service.UserHeaderUpdate{}
			for k, v := range list.([]*model.CrPersonFound) {
				fmt.Println(k)
				go func(k int, v *model.CrPersonFound) {

					wg.Add(1)
					c <- 1 // 作用类似于waitgroup.Add(1)
					fileName := httpHelp.DownImage(pwdDir, v.UserHeader)
					userHeaderData.UserHeader = fileName
					fmt.Println(v.Id)
					foundPersonService.UpdateDateById(v.Id, userHeaderData)
					defer wg.Done()
					<-c // 执行完毕，释放资源
					fmt.Println(k)

				}(k, v)
			}
			page++
			wg.Wait()

		}
	},
}

// getDirPath
//  @Description:
//  @receiver f
//  @param path
//  @param dir
//  @return string
func (f *FixCommand) getDirPath(path string, dir string) string {
	var sysType string
	if sysType == "windows" {
		// LINUX系统
		return path + "\\source\\image\\" + dir + "\\"

	} else {
		return path + "/source/image/" + dir + "/"

	}
}
func init() {
	rootCmd.AddCommand(FixImageCommand)
	FixImageCommand.Flags().String("dir", "开始了", "随便写的，我开心")
	FixImageCommand.Flags().String("web", "hopexr", "随便写的，我开心")
}
