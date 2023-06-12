package help

import (
	"api_client/constant"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var Appini string

type Common struct {
}

// 设置系统默认值
func (c *Common) SetDefault(typs reflect.Type, values reflect.Value) {

	for i := 0; i < values.NumField(); i++ {
		t := typs.Field(i)
		v := values.Field(i)

		if t.Name[0] >= 'a' && t.Name[0] <= 'z' {
			continue
		}

		var kind = t.Type.Kind()
		if kind == reflect.Struct {
			c.SetDefault(t.Type, v)
			continue
		} else if kind.String() == "ptr" {
			item := v.Interface()
			c.SetDefault(reflect.TypeOf(item).Elem(), reflect.ValueOf(item).Elem())
			continue
		}

		val := t.Tag.Get("default")
		if val == "" {
			continue
		}
		switch kind {
		case reflect.Bool:
			if v.Bool() == false && val == "true" {
				v.SetBool(true)
			}
		case reflect.String:
			if v.String() == "" {
				v.SetString(val)
			}
			//else{
			//	// 获取defaulut数据
			//	matched := regexp.MustCompile("([a-z]+)")
			//	constString := matched.FindAllString(val, -1)
			//
			//}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if v.Int() == 0 {
				v.SetInt(c.ToInt(val))
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if v.Uint() == 0 {
				v.SetUint(c.ToUint(val))
			}
		}
	}
}

/**
根据时间获取时间格式
*/
func (c *Common) GetTimeStamp(timeString string, format string) int {

	// 匹配成需要的格式
	timeSlice := strings.Split(timeString, "-")
	newtimesString := ""
	for _, l := range timeSlice {
		timesString := fmt.Sprintf("%02s", l)
		newtimesString = newtimesString + "-" + timesString
	}
	// 删除第一个字符 -
	timeString = strings.Replace(newtimesString, "-", "", 1)

	loc, _ := time.LoadLocation("Asia/Shanghai")

	tt, err := time.ParseInLocation(format, timeString, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
	if err != nil {
		panic("时间格式异常" + err.Error())
	}
	timeStamp := tt.Unix()
	// 转成字符串
	strInt64 := strconv.FormatInt(timeStamp, 10)
	timeStampInt, _ := strconv.Atoi(strInt64)

	return timeStampInt
}

//
func (c *Common) ToInt(val string) int64 {
	Num, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		panic(err.Error())
	}
	return Num
}

// 字符串uint64
func (c *Common) ToUint(val string) uint64 {
	u64, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		panic(err.Error())
	}
	return u64
}

// 接口转字符
func (c *Common) InterFaceToString(value interface{}) string {
	// Strval 获取变量的字符串值
	// 浮点型 3.0将会转换成字符串3, "3"
	// 非数值或字符类型的变量将会被转换成JSON格式字符串
	// interface 转 string
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}

// 判断map 是否在切片中
func isValueInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false

}

// CheckMobile 检验手机号
func (c *Common) CheckMobile(phone string) bool {
	// 匹配规则
	// ^1第一位为一
	// [345789]{1} 后接一位345789 的数字
	// \\d \d的转义 表示数字 {9} 接9位
	// $ 结束符
	regRuler := "^1[345789]{1}\\d{9}$"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(phone)

}

// CheckIdCard 检验身份证
func (c *Common) CheckIdCard(card string) bool {
	//18位身份证 ^(\d{17})([0-9]|X)$
	// 匹配规则
	// (^\d{15}$) 15位身份证
	// (^\d{18}$) 18位身份证
	// (^\d{17}(\d|X|x)$) 18位身份证 最后一位为X的用户
	regRuler := "(^\\d{15}$)|(^\\d{18}$)|(^\\d{17}(\\d|X|x)$)"

	// 正则调用规则
	reg := regexp.MustCompile(regRuler)

	// 返回 MatchString 是否匹配
	return reg.MatchString(card)
}

// 提取字符串中的手机号
func (c *Common) StrPhoneRule(str string) string {
	var param string
	Regexp := regexp.MustCompile(`([\d]{11})`)
	params := Regexp.FindStringSubmatch(str)
	if params == nil {
		return ""
	}
	param = params[0]
	return param
}

func (c *Common) substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func (c *Common) GetParentDirectory(dirctory string) string {
	sysType := runtime.GOOS
	if sysType == "windows" {
		return c.substr(dirctory, 0, strings.LastIndex(dirctory, "\\"))
		// LINUX系统
	} else {
		return c.substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
	}
}

func (c *Common) GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// GetFormatTime
//  @Description:
//  @receiver c
//  @param timeString 时间字符串
//  @return string
//  @date ${DATE} ${TIME}
func (c *Common) GetFormatTime(timeString string) string {
	if timeString != "0" {
		pushTimeInt64, _ := strconv.ParseInt(timeString, 10, 64)
		timeString = time.Unix(pushTimeInt64, 0).Format("2006-01-02")
	} else {
		timeString = ""
	}
	return timeString
}

// getImageUrl
//  @Description:  获取图片url
//  @receiver c
func (c *Common) GetImageUrl(imgUrl string) string {

	pos := strings.Index(imgUrl, "http")
	// 如果没有匹配到
	if pos < 0 {
		var imgPath string
		imgPath = constant.ImageUrl + "image/xunren/" + imgUrl
		return imgPath
	} else {
		return imgUrl
	}

}
