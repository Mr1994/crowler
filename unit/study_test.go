package unit

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"unicode/utf8"
)

func TestCalc(t *testing.T) {

	fmt.Println(utf8.RuneCountInString("你23好啊"))
	//name := utf8.RuneCountInString("你好啊")
	nameRune := []rune("你23好啊")
	fmt.Println("string(nameRune[:4]) = ", string(nameRune[:2]))
}

func TestType(t *testing.T) {
	var a *int    // 存储的是int的指针，目前为空
	var b int = 4 // 存储的是int的值
	a = &b        // a 指向 b 的地址
	//// a = b // a 无法等于 b，会报错，a是指针，b是值，存储的类型不同
	//fmt.Println(a) // a:0xc00000a090(返回了地址)
	//fmt.Println(*a) // *a:4(返回了值)
	//fmt.Println(&*a) // *抵消了&，返回了0xc00000a090本身
	*a = 5 // 改变 a 的地址的值
	fmt.Println(b)
	fmt.Println(*&b) // b:5，改变后 b 同样受到改变，因为 a 的地址是指向 b 的

}

type y struct {
	sync.RWMutex
	m map[int]int
}

func TestMap(t *testing.T) {

	// 定义结构体并赋值
	var m sync.Map
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key-%d", i)
		val := i + 1
		go func() {
			m.Store(key, val)
		}()
	}
	//time.Sleep(time.Second * 1)

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key-%d", i)
		x, _ := m.Load(key)
		fmt.Println(x)
	}
}

func TestChannel(t *testing.T) {

}

func TestJson(t *testing.T) {
	s := `{\"text\":[\"特征描述： \",\"婴儿时被放在苏州市东北街到娄门一段，后被送至姑苏孤儿院。本姓可能姓姚。\",\"失散经过： \",\"姓姚\",\"家属附言： \"]}`
	jsonStr, _ := strconv.Unquote("\"" + s + "\"") //直接把json字符串用Unquote处理一次即可

	fmt.Println(jsonStr)
	type MissDetail struct {
		Text string `json:"text"`
		Img  string `json:"img"`
	}
	var missDetail *MissDetail
	err := json.Unmarshal([]byte(s), &missDetail)
	fmt.Println(err.Error())
	fmt.Println(missDetail)

}
