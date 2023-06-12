package constant

//send_type类别
const SendTypeAuto = 0       //自动选择
const SendTypeXiaomi = 1     // 小米
const SendTypeHuawei = 2     // 华为
const SendTypeJiguang = 3    // 免费极光
const SendTypeVivo = 4       // vivo push
const SendTypeOppo = 5       // oppo push
const SendTypeJiguangVip = 7 // 免费极光
const SendTypeApple = 8      // 苹果apns

const PriorityLevel1 = 1

const PriorityLevel2 = 2

const GroupPriorityLevel1 = "group_JULIVE_QUEUE_PUSH_SERVICE_high_priority"
const GroupPriorityLevel2 = "group_JULIVE_QUEUE_PUSH_SERVICE_low_priority"

// 设置优先级
var PriorityMapGroupId = map[int]string{
	PriorityLevel1: GroupPriorityLevel1,
	PriorityLevel2: GroupPriorityLevel2,
}
