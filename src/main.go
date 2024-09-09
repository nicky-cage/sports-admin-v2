package main

import (
	"flag"
	"fmt"
	"os"
	"sports-admin/caches"
	"sports-admin/libs"
	"sports-admin/middlewares"
	"sports-admin/router"
	"sports-common/config"
	"sports-common/consts"
	"sports-common/log"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	//_ "net/http/pprof" // -- 检查性能问题时使用
)

// 默认定义的一些参数
var (
	versionName        = "sports-admin" //版本名称
	versionNameStr     = "1.0.0"        //版本号
	help               bool             //是否显示帮助信息
	isShowVersion      bool             //是否显示版本号
	isShowDesc         bool             //显示该版本的描述信息
	disableSysMaintain bool             //默认不禁用系统维护的标记
	configFilePath     string           //主配置
	extConfigFilePath  string           //微服务自身额外配置
	versionDesc        = "版本描述"         //版本描述
)

// 初始化相关操作
func init() {
	flag.BoolVar(&help, "h", false, "help info")
	flag.BoolVar(&isShowVersion, "v", false, "version 显示代码版本号")
	flag.BoolVar(&isShowDesc, "vd", false, "该版本的描述信息")
	flag.BoolVar(&disableSysMaintain, "m", false, "false系统维护的标记起作用 true表示后台维护的标记有不起作用")
	flag.StringVar(&configFilePath, "config", "", `必填，setting.ini配置文件的绝对路径,如"D:\shares\go-v\src\integrated-game-api\bin\setting.ini"`)
	flag.StringVar(&extConfigFilePath, "ext", "", `必填，setting_ext.ini配置文件的绝对路径,如"D:\shares\go-v\src\integrated-game-api\bin\setting_ext.ini"`)
	flag.Usage = func() { // 重写flag的usage,如果重写，默认的是使用flag.Boovar等中设置帮助信息说明
		_, _ = fmt.Fprintf(os.Stderr, `admin version: 1.13 beta
		Usage: -config=configvalue -ext=extConfigFile
		`)
		flag.PrintDefaults()
	}
}

// 加载判断启动参数
func parseParameters() {
	if help {
		flag.Usage()
		os.Exit(0)
	} else if isShowVersion {
		fmt.Println(versionName + versionNameStr)
		os.Exit(0)
	} else if isShowDesc {
		fmt.Println(versionName + versionNameStr + "\n" + versionDesc)
		os.Exit(0)
	}
}

// main
func main() {

	// 重新定义版本号
	buildVersion := time.Now().Format("20060102_1504")
	versionNameStr += "_build_" + buildVersion

	// 初始加载参数
	flag.Parse()
	parseParameters()
	fmt.Println(versionName + versionNameStr + "\n" + versionDesc)
	libs.LoadSelfExtConfigFile(extConfigFilePath)

	func() { // 设置常用需要加载的变量
		consts.CustomDebug = libs.IniGet("platform.custom_debug") == "1" //自定义debug是否开启
		consts.DisableSysMaintain = disableSysMaintain
		consts.AppName = libs.IniGet("app_name")
	}()

	// 初始化自身微服务的全局变量
	config.LoadConfigs(configFilePath)
	log.Start() //启动日志程序

	// pprof - 检查性能问题时使用
	// go func() {
	// 	log.Logger.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	// 启动主程序
	func() {
		caches.Initialize()  // 加载缓存
		app := gin.Default() // Gin
		// app.Use(middlewares.LogResponseBody)
		app.Use(middlewares.Logger()) // 日志
		router.Init(app)              // 路由
		_ = app.Run(libs.IniGet("service.port"))
	}()
}
