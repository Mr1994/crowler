package config

import (
	"database/sql"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DBConfig map[string]map[string]string // 数据库配置
var DB *gorm.DB                           // 业务库
var Hr *sql.DB                            // 人事库
var KafkaSyncClient sarama.SyncProducer   // kafka地址
var KafkaAsyncClient sarama.AsyncProducer // kafka地址
var KafkaConsumerClient sarama.Consumer   // kafka地址

//var ROOT_PATH string
var RDS *redis.Pool   // redis初始化
var IniConf *ini.File // ini路径

func Init() {

	//ROOT_PATH, _ = filepath.Abs(filepath.Dir(path) + "/../")
	DBConfig = map[string]map[string]string{
		"default": {
			"dialect":      IniConf.Section("DB").Key("Dialect").String(),
			"dsn":          IniConf.Section("DB").Key("DSN").String(),
			"maxIdleConns": IniConf.Section("DB").Key("MAX_IDLE_CONN").String(),
			"maxOpenConns": IniConf.Section("DB").Key("MAX_OPEN_CONN").String(),
		},
		"hr": {
			"dialect":      IniConf.Section("HR").Key("Dialect").String(),
			"dsn":          IniConf.Section("HR").Key("DSN").String(),
			"maxIdleConns": IniConf.Section("HR").Key("MAX_IDLE_CONN").String(),
			"maxOpenConns": IniConf.Section("HR").Key("MAX_OPEN_CONN").String(),
		},
	}
	SetDB()
	//SetHrDB()
	//SetRedisDb()
	//SetSyncKafka()
	//SetAsyncKafka()
}

//// 获取默认db
func SetDB() *gorm.DB {
	var err error
	DB, err = gorm.Open(mysql.Open(DBConfig["default"]["dsn"]), &gorm.Config{})
	if err != nil {
		log.Panicln("err:", err.Error())
	}
	return DB
}

// 获取user db
func SetHrDB() *gorm.DB {
	var err error
	db, err := gorm.Open(mysql.Open(DBConfig["hr"]["dsn"]))
	if err != nil {
		log.Panicln("err:", err.Error())
	}

	return db
}

// 获取db连接
func GetDB(name string) *gorm.DB {
	var err error
	db, err := gorm.Open(mysql.Open(DBConfig[name]["DSN"]))
	if err != nil {
		log.Panicln("err:", err.Error())
	}
	return db
}

func SetRedisDb() {
	Conf := IniConf
	redisAddress := Conf.Section("REDIS").Key("Redis").String()
	passWord := Conf.Section("REDIS").Key("Password").String()
	RDS = &redis.Pool{
		MaxIdle:   3, /*最大的空闲连接数*/
		MaxActive: 8, /*最大的激活连接数*/
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisAddress, redis.DialPassword(passWord))
			if err != nil {
				return nil, err
			}
			return c, nil
		},
	}
}

func SetSyncKafka() {
	// 新建一个arama配置实例
	config := sarama.NewConfig()

	// WaitForAll waits for all in-sync replicas to commit before responding.
	config.Producer.RequiredAcks = sarama.WaitForAll

	// NewRandomPartitioner returns a Partitioner which chooses a random partition each time.
	config.Producer.Partitioner = sarama.NewRandomPartitioner

	config.Producer.Return.Successes = true
	kafkaAddress := IniConf.Section("KAFKA").Key("Kafka").String()
	// 新建一个同步生产者
	var err error
	KafkaSyncClient, err = sarama.NewSyncProducer([]string{kafkaAddress}, config)
	if err != nil {
		log.Panicln("err:", err.Error())
		return
	}
}
func SetAsyncKafka() {
	config := sarama.NewConfig()
	//等待服务器所有副本都保存成功后的响应
	config.Producer.RequiredAcks = sarama.WaitForAll
	//随机向partition发送消息
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	//设置使用的kafka版本,如果低于V0_10_0_0版本,消息中的timestrap没有作用.需要消费和生产同时配置
	//注意，版本设置不对的话，kafka会返回很奇怪的错误，并且无法成功发送消息
	config.Version = sarama.V0_10_0_1
	//使用配置,新建一个异步生产者
	kafkaAddress := IniConf.Section("KAFKA").Key("Kafka").String()

	var err error
	KafkaAsyncClient, err = sarama.NewAsyncProducer([]string{kafkaAddress}, config)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func SetConsumer() {
	configConsumer := sarama.NewConfig()
	configConsumer.Consumer.Return.Errors = true
	configConsumer.Version = sarama.V0_11_0_2

	kafkaAddress := IniConf.Section("KAFKA").Key("Kafka").String()
	// consumer
	var err error
	KafkaConsumerClient, err = sarama.NewConsumer([]string{kafkaAddress}, configConsumer)

	if err != nil {
		panic("consumer_test create consumer error :" + err.Error())
		return
	}

}
