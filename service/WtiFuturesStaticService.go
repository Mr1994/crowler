package service

import (
	"api_client/help"
	"api_client/model"
	"encoding/json"
	"fmt"
	"regexp"
)

// 根据天数统计期货价格

type WtiFuturesStaticService struct {
}

// 更新期货数据
func (f *WtiFuturesStaticService) SaveFuturesStatic(futures model.WtiFutures) {

	// 获取
	var symbol = map[string]string{
		"symbol": futures.SerialNumber,
	}

	futuresUrlByDay := (&WtiFuturesStaticService{}).getFuturesUrlByDay(futures.SerialNumber)
	futuresStr := help.HttpGet(futuresUrlByDay, symbol)
	r := regexp.MustCompile("\\[.*\\]")
	futuresData := r.FindAllString(futuresStr, -1)

	var paramList []map[string]string
	json.Unmarshal([]byte(futuresData[0]), &paramList)
	fmt.Println(paramList)
	//var resMap []model.WtiFuturesStatic
	// 最大连接数10
	ch := make(chan int, 1)
	for _, v := range paramList {
		temp := &model.WtiFuturesStatic{}
		temp.Date = v["d"]
		temp.MaxPrice = v["h"]
		temp.OpenPrice = v["o"]
		temp.SerialNumber = futures.SerialNumber
		ch <- 1
		go temp.Create(temp, ch)
	}
}

func (f *WtiFuturesStaticService) getFuturesUrlByDay(serialNumber string) string {
	return "https://stock2.finance.sina.com.cn/futures/api/jsonp.php/var_" + serialNumber + "=/InnerFuturesNewService.getDailyKLine"
}
