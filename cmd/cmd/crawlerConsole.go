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
	"api_client/constant"
	"api_client/help"
	"api_client/model"
	"api_client/service"
	"encoding/json"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

var wgt sync.WaitGroup
var hasMore bool = true
var All string = "true"

// pushCmd represents the push command
var crawlerConsoleCmd = &cobra.Command{
	Use:   "crawler-console",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 确认是否是全量更新
		All, _ = cmd.Flags().GetString("all")
		if All == "" {
			panic("参数异常")
		}
		var crawerPageMap map[string]int
		// 初始化数据
		crawerPageMap = make(map[string]int)
		// 赋值 每个类别有多少页
		crawerPageMap = map[string]int{
			"http://www.xunrenla.com/xr/list_1_%d.html":             10,
			"http://www.xunrenla.com/jiarenxunqin/list_3_%d.html":   5,
			"http://www.xunrenla.com/qinrenxunjia/list_4_%d.html":   3,
			"http://www.xunrenla.com/xunzhaopengyou/list_6_%d.html": 4,
			"http://www.xunrenla.com/ganenxunren/":                  1,
			"http://www.xunrenla.com/taihaixunren/":                 1,
			"http://www.xunrenla.com/qita/":                         1,
		}
		for urlHttp, pageMax := range crawerPageMap {
			// 从第一页开始获取数据
			var page int = 1
			for page <= pageMax {
				var url string
				//ch := make(chan int, 1)
				if pageMax != 1 {
					url = fmt.Sprintf(urlHttp, page)
				} else {
					url = urlHttp
				}
				fmt.Println(url)
				docDetail := getHtmlDoc(url)
				// 获取list page里面的呢日用
				getPageListDoc(docDetail)
				if !hasMore {
					fmt.Println("进入下一次" + url)
					break
				}
				// 下一页
				page++
				if page == pageMax {
					continue
				}

			}
		}

	},
}

func getPageListDoc(docDetail *goquery.Document) {
	defer func() {
		err := recover()
		if err != nil {
			hasMore = false
		} else {
			hasMore = true
		}

	}()
	// docDetail 当前list html实体
	docDetail.Find(".siteImg a").Each(func(i int, selection *goquery.Selection) {
		//ch <- 1
		href, _ := selection.Attr("href")
		fmt.Println(href)
		docDetailDocument := getHtmlDoc(href)
		time.Sleep(1 * time.Second)
		// 定义指针结构体，并赋值空
		var cpPersonStruct *model.CrPerson
		// 初始化当前结构体数据
		cpPersonStruct = new(model.CrPerson)
		//crawlerStruct := service.FoundPersonService{} // 赋值一个空结构体
		crawlerStruct := new(service.FoundPersonService) // 初始化一个空结构体

		cpPersonStruct.CrawlerAddress = href
		cpPersonStruct.Source = constant.SourceXunRen
		// 开始匹配
		docDetailDocument.Find(".site-info-img img").Each(func(i int, selectionDetail *goquery.Selection) {
			detailHeaderImg, _ := selectionDetail.Attr("src")
			cpPersonStruct.UserHeader = detailHeaderImg
		})
		docDetailDocument.Find(".layui-card-header em").Each(func(i int, selectionDetail *goquery.Selection) {
			publishTime := selectionDetail.Text()
			publishTime = strings.TrimSpace(publishTime)
			crawlerStruct.GetCrawlerPushTime(publishTime, cpPersonStruct)

		})
		docDetailDocument.Find(".layui-card-header h2 span").Each(func(i int, selectionDetail *goquery.Selection) {
			crawlerId := selectionDetail.Text()
			strIndex := strings.Index(crawlerId, "ID")
			// 如果存在姓名字符，则添加到结构体中
			if strIndex >= 0 {
				strPhIndex := strings.LastIndex(crawlerId, "：")
				crawlerId = crawlerId[strPhIndex+3 : len(crawlerId)-3]
				//crawlerIdInt, _ := strconv.Atoi(crawlerId)
				cpPersonStruct.CrawlerId = crawlerId
			}

		})
		// 获取详情描述
		textSlice := []string{}
		docDetailDocument.Find(".article-content .layui-col-md12:nth-child(1) p").Each(func(i int, selectionDetail *goquery.Selection) {
			detailText := selectionDetail.Text()
			detailText = strings.TrimSpace(detailText)
			if detailText != "" {
				textSlice = append(textSlice, detailText)
			}
		})
		// 获取详情内的图片
		imgSlice := []string{}
		docDetailDocument.Find(".article-content .layui-col-md12:nth-child(1) img").Each(func(i int, selectionDetail *goquery.Selection) {
			detailImg, _ := selectionDetail.Attr("src")
			detailImg = "http://www.xunrenla.com" + detailImg
			imgSlice = append(imgSlice, detailImg)
		})
		var MissDetail = map[string][]string{}
		MissDetail["img"] = imgSlice
		MissDetail["text"] = textSlice
		missDetailJson, _ := json.Marshal(MissDetail)
		cpPersonStruct.MissDetail = string(missDetailJson)

		// 获取手机号码的内容
		docDetailDocument.Find(".site-info-list li a").Each(func(i int, selectionDetail *goquery.Selection) {
			detailText, _ := selectionDetail.Attr("href")
			crawlerStruct.AssemblyCrawlerData(detailText, cpPersonStruct)

		})
		// 获取文档内容文本的
		docDetailDocument.Find(".site-info-list li").Each(func(i int, selectionDetail *goquery.Selection) {
			detailText := selectionDetail.Text()
			crawlerStruct.AssemblyCrawlerData(detailText, cpPersonStruct)

		})

		//<-ch

		//close(ch)
		CrPersonModel := model.CrPerson{}
		res := CrPersonModel.Create(cpPersonStruct, All)
		// 如果不是全部的数据，则发现最新的就停止
		if !res && All == "false" {
			panic("数据已存在，可以停止了")
		}
	})

}

// pushCmd represents the push command
var crawlerConsoleHopeXrCmd = &cobra.Command{
	Use:   "crawler-console-hope-xr",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		All, _ = cmd.Flags().GetString("all")
		if All == "" {
			panic("参数异常")
		}
		fmt.Println("开始抓取")
		var crawerPageMap = []string{
			"https://www.manlost.com/search.html",
		}

		for _, urlHttp := range crawerPageMap {
			var i int
			i = 1
			//ch := make(chan int, 1)
			for {
				//ch <- 1
				var url string
				//// 从第一页开始获取数据
				url = fmt.Sprintf(urlHttp)
				fmt.Println(url)

				var docListHtml *goquery.Document
				docListHtml = new(goquery.Document)
				var pageMap = map[string]string{
					"page": strconv.Itoa(i),
				}

				docListHtml = getHtmlDocByParams(url, pageMap)

				r := docListHtml.Find(".our-faq")
				if len(r.Nodes) == 0 {
					fmt.Println("已经抓取到最后一页了")
					break
				}
				fmt.Println(r)
				// 写入一条数据
				//go getHopeXr(r, ch)
				i++

			}
		}

	},
}

func getHopeXr(r *goquery.Document) {

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
	r.Find(" .fj_post .details .job_chedule  +a").Each(func(i int, selection *goquery.Selection) {

		detailUrl, _ := selection.Attr("href")
		fmt.Println("开始了详情页抓取")

		fmt.Println(detailUrl)
		fmt.Println("开始了详情页抓取")
		docDetailHtml := getHtmlDoc(detailUrl)
		// 定义指针结构体，并赋值空
		var cpPersonStruct *model.CrPerson
		// 初始化当前结构体数据
		cpPersonStruct = new(model.CrPerson)
		crawlerStruct := new(service.FoundPersonService) // 初始化一个空结构体
		cpPersonStruct.Source = constant.SourceHopeXunRen

		// 抓取详细地址
		cpPersonStruct.CrawlerAddress = detailUrl
		// 名字
		cpPersonStruct.Name = strings.TrimSpace(selection.Text())
		// 丢失详情
		var textSlice []string
		textSlice = []string{}
		docDetailHtml.Find(".candidate_about_info p").Each(func(i int, selection *goquery.Selection) {
			text := selection.Text()
			trimString := strings.TrimSpace(text)
			regExp, _ := regexp.Compile(" ")
			rep2 := regExp.ReplaceAllString(text, "")
			if trimString != "" {
				textSlice = append(textSlice, rep2)
			}
		})

		var MissDetail map[string][]string
		MissDetail = make(map[string][]string)
		MissDetail["text"] = textSlice
		missDetailJson, _ := json.Marshal(MissDetail)
		cpPersonStruct.MissDetail = string(missDetailJson)

		// 类型
		cpPersonStruct.Type = 2

		// 只要下面第一级div内的文本
		selection.Find("img").Each(func(i int, selectionimg *goquery.Selection) {
			img, _ := selectionimg.Attr("src")
			cpPersonStruct.UserHeader = img
		})
		// 获取抓取编号
		crawlerStruct.GetCrawlerId(detailUrl, cpPersonStruct)

		// 抓取其他
		common := help.Common{}

		docDetailHtml.Find(".candidate_working_widget").Each(func(i int, selectiondiv *goquery.Selection) {
			// 时间
			selectiondiv.Find(" div:nth-child(2) > p").Each(func(i int, selectionTime *goquery.Selection) {
				matchStr := selectionTime.Text()
				fmt.Println("失踪时间")
				fmt.Println(matchStr)
				fmt.Println("失踪时间")

				missDay := common.GetTimeStamp(matchStr, "2006-01-02")
				cpPersonStruct.MissDay = missDay
			})
			selectiondiv.Find("div:nth-child(8) > p").Each(func(i int, selectionsex *goquery.Selection) {
				sex := strings.Replace(selectionsex.Text(), "性别：", "", 1)

				var sexInt int
				if sex == "男" {
					sexInt = constant.Man
				} else {
					sexInt = constant.WoMan

				}
				cpPersonStruct.Sex = sexInt
				fmt.Println(cpPersonStruct.Sex)
			})
		})

		CrPersonModel := model.CrPerson{}
		if All == "false" {
			pushTimeInt64 := time.Now().Unix()
			strInt64 := strconv.FormatInt(pushTimeInt64, 10)
			pushTimeInt, _ := strconv.Atoi(strInt64)
			cpPersonStruct.PushTime = pushTimeInt
		}
		res := CrPersonModel.Create(cpPersonStruct, All)
		if !res && All == "false" {
			panic("数据已存在，可以停止了")
		}

	})
}

// pushCmd represents the push command

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		println("开始抓取")
		var url string
		path := "cmd/cmd/cccbakx.xlsx"
		f, e := excelize.OpenFile(path)

		if e != nil {
			fmt.Println(e)
		}

		//f.InsertCol("apps", "H")
		//f.InsertCol("apps", "H")
		//os.Exit(1)

		sheetName := "cccbak"
		rows := f.GetRows(sheetName)

		for k, row := range rows {
			code := row[0]
			fmt.Printf("%06s", code)
			fmt.Println(k)
			url = fmt.Sprintf("http://stockpage.10jqka.com.cn/%s/", code)

			docListHtml := getHtmlDoc(url)
			time.Sleep(1 * time.Second)
			docListHtml.Find(".company_details dd:nth-child(4) ").Each(func(i int, selection *goquery.Selection) {
				text, _ := selection.Attr("title")
				fmt.Println(text)
				fmt.Println(k)
				f.SetCellValue(sheetName, fmt.Sprintf("BH%d", k+1), text)
			})
		}
		if err := f.SaveAs(path); err != nil {
			fmt.Println(err)
		}

	},
}

func getHtmlDoc(url string) *goquery.Document {
	var params = make(map[string]string)
	//
	body := help.HttpGet(url, params)
	isUtf := utf8.Valid([]byte(body))
	var readerBody io.Reader
	readerBody = strings.NewReader(body)
	// 如果不是uft-8 则转化成uft-8
	if !isUtf {
		readerBody = transform.NewReader(readerBody, simplifiedchinese.GBK.NewDecoder())
	}

	docDetail, err := goquery.NewDocumentFromReader(readerBody)
	if err != nil {
		panic(err)
	}
	return docDetail
}

func getHtmlDocByParams(url string, params map[string]string) *goquery.Document {
	//
	body := help.HttpGet(url, params)
	isUtf := utf8.Valid([]byte(body))
	var readerBody io.Reader
	readerBody = strings.NewReader(body)
	// 如果不是uft-8 则转化成uft-8
	if !isUtf {
		readerBody = transform.NewReader(readerBody, simplifiedchinese.GBK.NewDecoder())
	}

	docDetail, err := goquery.NewDocumentFromReader(readerBody)
	if err != nil {
		panic(err)
	}
	return docDetail
}

// pushCmd represents the push command
var crawlerConsoleHopeXrNewCmd = &cobra.Command{
	Use:   "crawler-console-hope-xr-new",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:
	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		All, _ = cmd.Flags().GetString("all")
		if All == "" {
			panic("参数异常")
		}
		fmt.Println("开始抓取")
		urlHttp := "https://www.manlost.com/search.html"

		var i int
		i = 1
		//ch := make(chan int, 1)
		for {

			var docListHtml *goquery.Document
			docListHtml = new(goquery.Document)
			var pageMap = map[string]string{
				"page": strconv.Itoa(i),
			}

			docListHtml = getHtmlDocByParams(urlHttp, pageMap)

			r := docListHtml.Find(".our-faq")
			if len(r.Nodes) == 0 {
				fmt.Println("已经抓取到最后一页了")
				break
			}
			// 写入一条数据
			//go getHopeXr(docListHtml, ch)
			getHopeXr(docListHtml)
			i++

		}

	},
}

func init() {
	rootCmd.AddCommand(crawlerConsoleCmd)          // 寻人网
	rootCmd.AddCommand(crawlerConsoleHopeXrNewCmd) // 希望寻人网
	rootCmd.AddCommand(crawlerConsoleHopeXrCmd)    // 希望寻人网
	rootCmd.AddCommand(testCmd)                    // 寻狗网
	crawlerConsoleCmd.Flags().String("all", "", "true 标识全量数据")
	crawlerConsoleHopeXrCmd.Flags().String("all", "", "true 标识全量数据")
	crawlerConsoleHopeXrNewCmd.Flags().String("all", "", "true 标识全量数据")
}
