package unit

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"
	"unsafe"
)

func TestSlice(t *testing.T) {
	// 数组
	var a = [2][2]int{{1, 1}, {2, 2}}
	fmt.Println(a)
	// 切片
	var s []int
	s = make([]int, 10)
	s = []int{2, 3, 5, 7, 11, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13, 13}
	printSlice(s)
	s = s[:0]
	printSlice(s)

	s = s[:4]
	printSlice(s)

	s = s[2:]
	printSlice(s)

	for i, n := 0, 10; i < n; i++ {
		s = append(s, 1)
		fmt.Printf("len=%d cap=%d\n", len(s), cap(s))
	}

}

func TestSlice2(t *testing.T) {
	a := [...]int{0, 1, 2, 3}
	x := a[:1]
	fmt.Println(x)
	y := a[2:]
	fmt.Println(y)

	x = append(x, y...)
	fmt.Println(a, x)

	x = append(x, y...)
	fmt.Println(a, x)
}

func TestRange(t *testing.T) {

	var m = []int{1, 2, 3}
	var wg sync.WaitGroup
	for i := range m {
		wg.Add(1)
		go func(i int) {
			fmt.Print(i)
			wg.Done()
		}(i)
	}
	wg.Wait()

	return
	slice := []int{0, 1, 2, 3}
	mp := make(map[int]*int)
	// 指向的是同一个地址，一个地址值改变， 会引起整个mp的index进行修改
	for index, value := range slice {
		mp[index] = &value
		fmt.Println("address is:", &value)
		fmt.Println("value is:", value)
	}
	//根本原因在于for range是用一个变量承接mp中的内容的
	fmt.Println("-------------------------------------------------------------------")
	for key, value := range mp {
		fmt.Println(key, " ", *value)
	}
	return
	mySlice := []string{"I", "am", "peachesTao"}
	for _, ele := range mySlice {
		ele = ele + "-new"
	}
	fmt.Println(mySlice)
}
func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}

// defer测试
func TestDefer(t *testing.T) {
	c := dfer(0)

	defer func() {
		fmt.Println("1号")
		fmt.Println(recover())
		if err := recover(); err != nil {
			fmt.Printf("异常了1号:%v\n", err)
		}

	}()
	defer func() {
		fmt.Println("2号")
		fmt.Println(recover())
		if err := recover(); err != nil {
			fmt.Printf("异常了2号:%v\n", err)
		}

	}()
	panic("11111")
	panic("2222")
	fmt.Println("回来的是")
	fmt.Println(*c)
}

func dfer(i int) *int {
	defer func() {
		i++
		fmt.Println("defer2:", i)

	}()
	defer func() {
		i++
		fmt.Println("defer1:", i)
	}()
	fmt.Println(i)
	return &i
}

func TestUint(t *testing.T) {
	var a uint8 = 1
	var b uint8 = 255
	fmt.Println("减法：", a-b)
	fmt.Println("加法：", a+b)
}

func TestRune(t *testing.T) {
	var str = "你好呀234啊"
	fmt.Println("截取string长度", str[:2])
	fmt.Println("获取string长度", len(str))

	var c = []byte(str)
	fmt.Println("截取byte长度", string(c[:2]))
	fmt.Println("获取byte长度", len(c))

	nameRune := []rune(str)
	fmt.Println("截取rune长度", string(nameRune[0:4]))
	fmt.Println("获取rune长度", len(nameRune))
}

func TestTag(t *testing.T) {
	type Person struct {
		Name string `json:"name" bson:"Name"`
		Age  int    `json:"age" bson:"Age"`
	}
	var person = Person{"nihao", 11}

	personType := reflect.TypeOf(person)
	fieldAge, isOk := personType.FieldByName("Age")
	fmt.Println(fieldAge)
	if isOk {
		jsonTag := fieldAge.Tag.Get("json")
		bsonTag := fieldAge.Tag.Get("bson")
		fmt.Println("Name Json Tag =", jsonTag, "Bson Tag =", bsonTag)
	} else {
		fmt.Println("No Age Field")
	}

	for i := 0; i < personType.NumField(); i++ {
		// 获取每个成员的结构体字段类型
		fieldType := personType.Field(i)
		// 输出成员名和tag
		fmt.Printf("name: %v  tag: '%v'\n", fieldType.Name, fieldType.Tag)
	}
	// 通过字段名, 找到字段类型信息
	if catType, ok := personType.FieldByName("Type"); ok {
		// 从tag中取出需要的tag
		fmt.Println(catType.Tag.Get("json"), catType.Tag.Get("id"))
	}

}

// TestPassValu
//  @Description: 传值测试
//  @param t
func TestPassValu(t *testing.T) {

	var cMap map[string]string
	cMap = make(map[string]string)

	cMap["af"] = "ddd"
	fmt.Printf("&第1次cMap value: %p\n ", cMap)
	fmt.Printf("&第1次cMap 地址: %p\n ", &cMap)
	cMap["a"] = "还好"
	fmt.Printf("&第2次cMap value: %p\n ", cMap)
	fmt.Printf("&第2次cMap 地址: %p\n ", &cMap)

	fmt.Println(&cMap)

	var cInt = 12
	fmt.Println(cInt)
	fmt.Println(&cInt)
	cInt = 13
	fmt.Println(cInt)
	fmt.Println(&cInt)

}

func TestSelect(t *testing.T) {
	fmt.Println("start")
	//ch := make(chan int, 1)
	//ch <- 1024
	//select {
	//case val := <-ch:
	//	fmt.Println("Received from ch1, val =", val)
	//case val := <-ch:
	//	fmt.Println("Received from ch2, val =", val)
	//case val := <-ch:
	//	fmt.Println("Received from ch3, val =", val)
	//default:
	//	fmt.Println("Run in default")
	//}
	/**
	符合哪一个进入哪一个========================
	*/
	//timeout := make(chan bool, 1)
	//ch := make(chan int,1)
	//ch <- 1024
	//
	//go func() {
	//	time.Sleep(1 * time.Second)
	//	timeout <- true
	//}()
	//select {
	//case <-ch:
	//	fmt.Println("received from ch")
	//case <-timeout:
	//	fmt.Println("select timeout")
	//}
	// 控制超时，防止死锁，没有default的时候========================================
	//ch := make(chan int)
	//select {
	//case <-ch:
	//	fmt.Println("received from ch")
	//case <-time.After(time.Second * 2):
	//	fmt.Println("select timer after timeout")
	//default:
	//	fmt.Println("11111111")
	//}
	// 循环中时候用select
	c1 := make(chan string)
	c2 := make(chan string)
	go func() {
		time.Sleep(time.Second * 1)
		c1 <- "one"
	}()
	go func() {
		time.Sleep(time.Second * 2)
		c2 <- "two"
	}()
	var i = 0
	o := make(chan bool)
	for {
		fmt.Println(i)
		select {
		case msg1 := <-c1:
			fmt.Println("received1", msg1)
		case msg2 := <-c2:
			fmt.Println("received2", msg2)
			o <- true
		case <-time.After(time.Second * 2):
			fmt.Println("没有了")
			o <- true
			goto ForEnd

		}
		i++
	}
ForEnd:
	<-o
	fmt.Println("程序结束")

}

func TestQuotation(t *testing.T) {
	//String in double quotes
	x := "tit\nfor\ttat"
	fmt.Println("Priting String in Double Quotes:")
	fmt.Printf("x is: %s\n", x)

	//String in back quotes
	y := `tit\nfor\ttat`
	fmt.Println("\nPriting String in Back Quotes:")
	fmt.Printf("y is: %s\n", y)

	//Declaring a byte with single quotes
	var b byte = 'a'
	fmt.Println("\nPriting Byte:")
	//Print Size, Type and Character
	fmt.Printf("Size: %d\nType: %s\nCharacter: %c\n", unsafe.Sizeof(b), reflect.TypeOf(b), b)

	//Declaring a rune with single quotes
	r := '£'
	fmt.Println("\nPriting Rune:")
	//Print Size, Type, CodePoint and Character
	fmt.Printf("Size: %d\nType: %s\nUnicode CodePoint: %U\nCharacter: %c\n", unsafe.Sizeof(r), reflect.TypeOf(r), r, r)
	//Below will raise a compiler error - invalid character literal (more than one character)
	//r = 'ab'

}

func TestMapBingFa(t *testing.T) {
	// 这种情况会报错
	//c := make(map[string]int)
	//go func() {
	//	for j := 0; j < 1000; j++ {
	//		c[fmt.Sprintf("%d", j)] = j
	//	}
	//}()
	//
	//go func() {
	//	for j := 0; j < 1000; j++ {
	//		fmt.Println(c[fmt.Sprintf("%d",j)])
	//	}
	//}()

	var c = struct {
		sync.RWMutex
		m map[string]int
	}{m: make(map[string]int)}
	wg := sync.WaitGroup{}

	wg.Add(1)
	//ch := make(chan struct{}, 1) // 控制最大并发数是 10
	go func() { //开一个goroutine写map
		for j := 0; j < 10000; j++ {
			c.Lock()
			//ch <- struct{}{}
			c.m[fmt.Sprintf("%d", j)] = j
			c.Unlock()
			//wg.Done()

		}
	}()

	go func() { //开一个goroutine读map
		for j := 0; j < 10000; j++ {
			c.RLock()
			fmt.Println(c.m[fmt.Sprintf("%d", j)])
			c.RUnlock()

		}
	}()
	wg.Done()

	wg.Wait()

	//time.Sleep(time.Second * 20)

}
