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
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	"strconv"
)

const AccessToken = "879ccb18e49b41cfe6dfb6a524b5cafcec8a7b4d9bdc654600ac297e74b54ab2"
const Screct = "SECa4b103d0ca5ac4e95a6858da40f0b42578db4f665ec73576faeab1ff8d9ef63d"

// pushCmd represents the push command
var crawlerWxConsoleCmd = &cobra.Command{
	Use:   "crawlerwx-console",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 确认是否是全量更新
		All, _ = cmd.Flags().GetString("all")
		url := "https://tophub.today/"
		html := getHtmlDoc(url)
		pushStr := getPmDoc(html)
		help.SendDingCommonMessage(pushStr, AccessToken, Screct)

	},
}

func getPmDoc(docDetail *goquery.Document) string {
	zhihuContent := ""
	//pmContent := ""
	caixinContent := ""
	keContent := ""

	// docDetail 当前list html实体
	maxNum := 5
	docDetail.Find(".c-d-e").Each(func(o int, selection *goquery.Selection) {

		docDetail.Find("#node-6 a").Each(func(i int, selection *goquery.Selection) {

			if i <= maxNum {
				s := i

				href, _ := selection.Attr("href")
				selection.Find(".t").Each(func(i int, selectionSon *goquery.Selection) {
					t := selectionSon.Text()
					zhihuContent += strconv.Itoa(s) + ": " + t + " " + " " + href + "\n\n"

				})

			}
		})
		zhihuContent = "【知乎TOP5 - 综合类】\n\n" + zhihuContent + "\n\n"

		//docDetail.Find("#node-221 a").Each(func(i int, selection *goquery.Selection) {
		//
		//	if i <= maxNum {
		//
		//		href, _ := selection.Attr("href")
		//		selection.Find(".t").Each(func(f int, selectionSon *goquery.Selection) {
		//			s := i
		//
		//			t := selectionSon.Text()
		//			pmContent += strconv.Itoa(s) + ": " + " " + t + " " + href + "\n\n"
		//
		//		})
		//
		//	}
		//})
		//pmContent = "【抖音TOP5 - 娱乐类】\n\n" + pmContent + "\n\n"

		docDetail.Find("#node-215 a").Each(func(i int, selection *goquery.Selection) {

			if i <= maxNum {

				href, _ := selection.Attr("href")
				selection.Find(".t").Each(func(h int, selectionSon *goquery.Selection) {
					s := i
					t := selectionSon.Text()
					caixinContent += strconv.Itoa(s) + ": " + " " + t + " " + href + "\n\n"

				})

			}
		})
		caixinContent = "【雪球TOP5 - 财经类】\n\n" + caixinContent

		docDetail.Find("#node-11 a").Each(func(i int, selection *goquery.Selection) {

			if i <= maxNum {

				href, _ := selection.Attr("href")
				selection.Find(".t").Each(func(h int, selectionSon *goquery.Selection) {
					s := i
					t := selectionSon.Text()
					keContent += strconv.Itoa(s) + ": " + " " + t + " " + href + "\n\n"

				})

			}
		})
		keContent = "【36氪TOP5 - 科技类】\n\n" + keContent

		return
	})

	pushStr := "【早间成长新闻】\n\n" + zhihuContent + caixinContent + keContent + "\n\n"

	return pushStr

}

func init() {
	rootCmd.AddCommand(crawlerWxConsoleCmd) // 初始化脚本
	crawlerWxConsoleCmd.Flags().String("all", "", "true 标识全量数据")
}
