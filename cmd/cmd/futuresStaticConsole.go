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
	"api_client/model"
	"api_client/service"
	"fmt"
	"github.com/spf13/cobra"
	"strconv"
)

// pushCmd represents the push command
var futuresStaticCmd = &cobra.Command{
	Use:   "futures-static",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//serialNumber, _ := cmd.Flags().GetString("serial_number")
		// 设置发送类型
		serialList := (&model.WtiFutures{}).Find()
		for _, v := range serialList {
			(&service.WtiFuturesStaticService{}).SaveFuturesStatic(v)
		}

	},
}

// pushCmd represents the push command
var futuresCalcCmd = &cobra.Command{
	Use:   "futures-calc",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		serialNumber, _ := cmd.Flags().GetString("serial_number")
		// 查询需要分析的当前所有关联数据
		serialList := (&model.WtiFuturesAssociate{}).Find(serialNumber)
		var data = make(map[string]map[string]map[string]string)
		var subData = make(map[string]map[string]string)

		for _, v := range serialList {
			// 查询当前的统计数据
			serialStaticList := (&model.WtiFuturesStatic{}).Find(v.SerialNumber)
			fmt.Println(serialStaticList)
			for _, v1 := range serialStaticList {
				var subSonData = make(map[string]string)
				subSonData["maxPrice"] = v1.MaxPrice
				subSonData["minPrice"] = v1.MinPrice
				subSonData["openPrice"] = v1.OpenPrice

				subData[v1.SerialNumber] = subSonData
				data[v1.Date] = subData
			}
		}
		//fmt.Println(data)
		// 数据对比 k 日期 v 编号->价格
		for k, v := range data {
			var str string = "当日"
			var successRatio = make(map[string]bool)
			for k1, v1 := range v {

				var number = k1
				maxPrice, _ := strconv.ParseFloat(v1["maxPrice"], 32)
				openPrice, _ := strconv.ParseFloat(v1["openPrice"], 32)
				// 计算同一天相似比例
				res := maxPrice - openPrice
				successRatio[number] = res > 0
				str += "日期" + k + ":品种" + k1 + "：最大价格" + v1["maxPrice"] + "：开盘价格" + v1["openPrice"]
			}
			fmt.Println(str)
			fmt.Println(successRatio)
		}

		//fmt.Println(data)
	},
}

func init() {
	rootCmd.AddCommand(futuresStaticCmd)
	rootCmd.AddCommand(futuresCalcCmd)
	futuresCalcCmd.Flags().String("serial_number", "", "期货编号")
}
