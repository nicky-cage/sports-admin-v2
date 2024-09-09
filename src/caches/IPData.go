package caches

import (
	"fmt"
	common "sports-common"
	"sports-common/consts"
	"sports-common/tools"
	"strings"
	"sync"

	"github.com/go-redis/redis/v7"
)

//
var ckIPList = "admin:ip_areas"

// IPData 结构
type IPData struct {
	Locker sync.Mutex
	List   map[string]string
}

// IPCached 默认
var IPCached = &IPData{
	List: map[string]string{},
}

// GetArea 得到ip对应的地区
func (ths *IPData) GetArea(ip string, autoAdd ...bool) string {
	isAutoAdd := false // 是否自动追加到缓存当中
	if len(autoAdd) >= 1 {
		isAutoAdd = autoAdd[0]
	}

	areas := strings.Trim(tools.GetAreaByIp(ip), "[]")
	areaArr := strings.Split(areas, ",")
	area := ""

	// 处理是否自动追加ip到缓存当中
	defer func() {
		if isAutoAdd { // 如果自动追加到redis当中
			if _, exists := ths.List[ip]; exists {
				return
			}

			ths.Locker.Lock()
			defer ths.Locker.Unlock()

			rd := common.Redis(consts.PlatformIntegrated)
			defer common.RedisRestore(consts.PlatformIntegrated, rd)
			_, _ = rd.HSet(ckIPList, ip, area).Result()
			ths.List[ip] = area
		}
	}()

	// 可能没有查到ip
	if len(areaArr) <= 1 {
		if areaArr[0] != "" {
			area = areaArr[0]
		} else {
			area = "*未知地区*"
		}
		return area
	}

	// 如果是本机地址
	if areaArr[0] == "本机地址" {
		area = "-本机地址-"
		return area
	}

	// 查到多个地址
	area += areaArr[0]
	if len(areaArr) <= 2 || areaArr[1] == "" {
		return area
	}
	area += "-" + areaArr[1]
	if len(areaArr) <= 3 || areaArr[2] == "" {
		return area
	}
	area += "-" + areaArr[2]
	return area
}

// Init 初始化
func (ths *IPData) Init(platform string) {
	rd := common.Redis(consts.PlatformIntegrated)
	defer common.RedisRestore(consts.PlatformIntegrated, rd)

	db := common.Mysql(platform)
	defer db.Close()

	sql := "SELECT last_login_ip, register_ip FROM users"
	rows, err := db.QueryString(sql)
	if err != nil {
		fmt.Println("获取IP(db)信息出错:", err)
		return
	}

	for _, r := range rows {
		lastLoginIP := r["last_login_ip"]
		registerIP := r["register_ip"]

		ths.AddIP(rd, lastLoginIP)
		ths.AddIP(rd, registerIP)
	}
}

// AddIP 将ip写入缓存当中
func (ths *IPData) AddIP(rd *redis.Conn, ip string) {
	realIP := strings.TrimSpace(ip)
	if realIP == "" {
		return
	}

	// 检查内存
	if _, exists := ths.List[ip]; exists { // 如果已经存在, 则不处理
		return
	}

	// 检查redis
	val, err := rd.HGet(ckIPList, ip).Result()
	if err != nil && err != redis.Nil {
		fmt.Println("获取ip(redis)信息出错:", err)
	}
	if val != "" { // 设置到内存当中
		ths.List[ip] = val
		return
	}

	area := ths.GetArea(ip)
	if area == "" {
		fmt.Println("获取ip发生错误:", ip)
		return
	}

	// 设置redis
	_, err = rd.HSet(ckIPList, ip, area).Result()
	if err != nil {
		fmt.Println("设置ip信息出错:", err)
		return
	}

	// 设置内存
	ths.List[ip] = area
}
