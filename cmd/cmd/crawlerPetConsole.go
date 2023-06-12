/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"api_client/help"
	"api_client/model"
	"api_client/service"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	urlNet "net/url"
	"os"
	"strings"
	"time"
)

// pushCmd represents the push command
var crawlerConsoleXunGouCmd = &cobra.Command{
	Use:   "crawler-console-xun-gou",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		println("开始抓取")
		// 确认是否是全量更新
		All, _ = cmd.Flags().GetString("all")
		if All == "" {
			panic("参数异常")
		}
		var i int = 1

		// 抓取当前page
		defer func() {
			err := recover()
			if err != nil {
				new(help.Common).Log(err.(string), "cmd")
				help.SendDingMessage(err.(string))
				if err.(string) == "数据已存在，可以停止了" {
					os.Exit(1)
				}
			}
		}()
		for {
			var url string
			//// 从第一页开始获取数据
			url = fmt.Sprintf("http://www.xungou5.com/category-catid-1-page-%d.html", i)
			fmt.Println(url)
			docListHtml := getHtmlDoc(url)
			r := docListHtml.Find(".media-body-title a")
			if len(r.Nodes) == 0 {
				fmt.Println("已经抓取到最后一页了")
				break
			}

			docListHtml.Find(".media-body-title a").Each(func(i int, selection *goquery.Selection) {
				href, _ := selection.Attr("href")
				fmt.Println(href)
				docDetailHtml := help.GetHtmlDoc(href)
				var crPetModel *model.CrPet
				crPetModel = new(model.CrPet)
				common := help.Common{}
				// 获取标题
				crPetModel.Source = 1
				crPetModel.CrawlerAddress = href
				// 解析url
				docDetailHtml.Find("#ajaxcomment script").Each(func(i int, selection *goquery.Selection) {
					urlNone, _ := selection.Attr("src")
					urlParams, e := urlNet.Parse(urlNone)
					if e != nil {
						panic(e.Error())
					}
					// 解析参数
					params, _ := urlNet.ParseQuery(urlParams.RawQuery)
					crPetModel.CrawlerId = params["id"][0]
				})
				docDetailHtml.Find(".information_title").Each(func(i int, selection *goquery.Selection) {
					text := selection.Text()
					crPetModel.Title = text
				})

				var MissText []string
				MissText = []string{}
				MissImg := []string{}
				MissDetail := map[string][]string{}
				var jsonTextDetail []byte
				// 拉去详情
				docDetailHtml.Find(".bd_left .maincon p ").Each(func(i int, selection *goquery.Selection) {
					// 匹配手机号码
					text := selection.Text()
					matchingPhone := common.StrPhoneRule(text)
					if matchingPhone != "" {
						crPetModel.ContractPhone = common.StrPhoneRule(text)
					}
					// 去掉换行符
					trimSpaceSlice := strings.Split(text, "\n")
					// 每行存储成一条数据
					for _, v := range trimSpaceSlice {
						trimSpace := strings.Replace(v, " ", "", -1)
						trimSpace = strings.Replace(v, "\t", "", -1)
						if trimSpace != "" {
							MissText = append(MissText, trimSpace)
						}
					}
					// 去掉空格
				})

				docDetailHtml.Find(".extra_contact li").Each(func(i int, selection *goquery.Selection) {
					text := selection.Text()
					service.AssemblyCrawlerPet(text, crPetModel)

				})
				// 抓取图片
				docDetailHtml.Find(".bd_left .bd img ").Each(func(i int, selection *goquery.Selection) {
					imgUrl, _ := selection.Attr("src")
					// 去掉换行符
					trimImgSlice := strings.Split(imgUrl, "\n")
					// 每行存储成一条数据
					for _, v := range trimImgSlice {
						MissImg = append(MissImg, v)
					}
					// 去掉空格
				})
				MissDetail["text"] = MissText
				MissDetail["img"] = MissImg
				jsonTextDetail, _ = json.Marshal(MissDetail)
				crPetModel.MissDetail = string(jsonTextDetail)
				res := crPetModel.Create(crPetModel, All)
				fmt.Println(res)
				if !res && All == "false" {
					panic("数据已存在，可以停止了")
				}
				time.Sleep(1 * time.Second)

			})

			i++
		}
	},
}

func init() {
	rootCmd.AddCommand(crawlerConsoleXunGouCmd) // 寻狗网
	crawlerConsoleXunGouCmd.Flags().String("all", "", "true 标识全量数据")

}
