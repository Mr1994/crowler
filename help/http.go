package help

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/wanghuiyt/ding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"unicode/utf8"
)

type HttpHelp struct {
}

// get方式http请求
func HttpGet(httpUrl string, paramsMap map[string]string) string {
	params := url.Values{}
	Url, err := url.Parse(httpUrl)
	if err != nil {
		panic("参数异常，请求url为空")
	}
	for k, v := range paramsMap {
		params.Set(k, v)
	}
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	fmt.Println(urlPath)

	resp, err := http.Get(urlPath)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp)
		panic("网络异常,请求地址" + urlPath)
	}
	return string(body)
}

// json 类型http请求
func HttpPost(httpurl string, paramsMap map[string]interface{}) string {

	mjson, _ := json.Marshal(paramsMap)
	requestBody := string(mjson)

	var jsonStr = []byte(requestBody)

	req, err := http.NewRequest("POST", httpurl, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)

	fmt.Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		panic("网络异常,请求地址" + err.Error())
	}
	return string(body)

}

func SendDingMessage(message string) {
	d := ding.Webhook{
		AccessToken: "879ccb18e49b41cfe6dfb6a524b5cafcec8a7b4d9bdc654600ac297e74b54ab2",    // 上面获取的 access_token
		Secret:      "SECa4b103d0ca5ac4e95a6858da40f0b42578db4f665ec73576faeab1ff8d9ef63d", // 上面获取的加签的值
	}
	x := d.SendMessage(message)
	fmt.Println(x)
}
func SendDingCommonMessage(message, AccessToken, Secret string) {

	d := ding.Webhook{
		AccessToken: AccessToken, // 上面获取的 access_token
		Secret:      Secret,      // 上面获取的加签的值
	}
	x := d.SendMessage(message)
	fmt.Println(x)
}

/**
获取html doc内容
*/
func GetHtmlDoc(url string) *goquery.Document {
	var params = make(map[string]string)
	//
	body := HttpGet(url, params)
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

func (h *HttpHelp) DownImage(imgPath string, imgUrl string) string {
	fileName := path.Base(imgUrl)
	res, err := http.Get(imgUrl)
	if err != nil {
		fmt.Println("A error occurred!")
		return imgPath
	}
	defer res.Body.Close()
	// 获得get请求响应的reader对象
	reader := bufio.NewReaderSize(res.Body, 32*1024)

	file, err := os.Create(imgPath + fileName)
	if err != nil {
		panic(err)
	}
	// 获得文件的writer对象
	writer := bufio.NewWriter(file)
	io.Copy(writer, reader)

	return fileName

}
