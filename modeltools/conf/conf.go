package conf

import "api_client/config"

// 是否覆盖已存在model
const ModelReplace = true

// 数据库驱动
const DriverName = "mysql"

type DbConf struct {
	Host   string
	Port   string
	User   string
	Pwd    string
	DbName string
}

// 数据库链接配置
var MasterDbConfig DbConf

// 获取创建model配置地址
func GetCreateModelConf() DbConf {
	MasterDbConfig.Host = config.IniConf.Section("CM").Key("Host").String()
	MasterDbConfig.Port = config.IniConf.Section("CM").Key("Port").String()
	MasterDbConfig.User = config.IniConf.Section("CM").Key("User").String()
	MasterDbConfig.Pwd = config.IniConf.Section("CM").Key("Password").String()
	MasterDbConfig.DbName = config.IniConf.Section("CM").Key("Dbname").String()
	return MasterDbConfig
}

func GetCmPath() string {
	return config.IniConf.Section("MPATH").Key("Mpath").String()
}
