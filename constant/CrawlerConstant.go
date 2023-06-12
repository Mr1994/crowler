package constant

const Man int = 1
const WoMan int = 2

var FindType = map[string]int{
	"家人寻亲": 1,
	"亲人寻家": 2,
	"感恩寻人": 3,
	"寻找朋友": 4,
	"寻找战友": 5,
	"寻找同学": 6,
	"台海寻亲": 7,
	"其他寻人": 8,
}

/**
希望寻人网状态映射
*/
var HopeXrFindType = map[string]int{
	"家寻亲人": 1,
	"亲人寻家": 2,
	"感恩寻人": 3,
	"寻找老友": 4,
	"战友情深": 5,
	"台海寻人": 7,
	"其他寻人": 8,
}

type Class struct {
	ClassNo   uint64
	ClassName string
	Students  string //数组
	IsHonor   bool
}

const SourceXunRen = 1     //寻人网 http://www.xunrenla.com
const SourceHopeXunRen = 2 //希望寻人网 https://www.hopexr.com

const CallPoliceYes = 1 // 报警
const CallPoliceNo = 2  // 未报警
