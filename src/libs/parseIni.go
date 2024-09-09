package libs

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

// selfIni 自身配置文件
var selfIni *ini.File

// LoadSelfExtConfigFile 加载自身配置文件
func LoadSelfExtConfigFile(configFile string) {
	if configFile == "" {
		configFile = "setting_ext.ini"
	}
	cfg, err := ini.Load(configFile) //加载配置文件
	if err != nil {
		fmt.Printf("读取扩展配置文件出错: %v\n", err)
		os.Exit(1)

	}
	selfIni = cfg //设置全局Ini信息
}

// GetIniKey 依据key获得ini的设置项
func GetIniKey(key string) *ini.Key {
	sArr := strings.Split(key, ".")
	arrLen := len(sArr)
	if arrLen == 1 {
		return selfIni.Section("").Key(key)
	} else if arrLen == 2 {
		return selfIni.Section(sArr[0]).Key(sArr[1])
	}
	return nil
}

// IniGet 获取配置项信息, 注意: 此方法必须在 LoadConfigs() 方法之后执行, 否则将无法获取配置信息
// key: 配置项名称
// 顶级配置项获取: config.Get("app_name")
// 二级配置项获取: config.Get("database.name")
func IniGet(key string) string {
	if selfIni == nil {
		return ""
	}
	return GetIniKey(key).String()
}
