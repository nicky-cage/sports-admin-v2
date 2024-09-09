package caches

import (
	"encoding/json"
	common "sports-common"
	"sports-common/consts"
	"sports-common/tools"

	"github.com/go-redis/redis/v7"
)

const CacheAdminPrefix = "admin:cache:v2:"

// getCache 从缓存当中获取对象
var getCache = func(platform, key string, obj interface{}) error {
	rd := common.Redis(platform)
	defer common.RedisRestore(platform, rd)
	val, err := rd.Get(CacheAdminPrefix + key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(val), obj)
	if err != nil {
		return err
	}

	return nil
}

// getCacheByJson 获取json源信息
var getCacheByJson = func(platform, key string) string {
	rd := common.Redis(platform)
	defer common.RedisRestore(platform, rd)
	val, err := rd.Get(CacheAdminPrefix + key).Result()
	if err != nil || err == redis.Nil {
		return ""
	}
	return val
}

// setCache 写入对象到缓存当中
var setCache = func(platform, key string, obj interface{}) error {
	rd := common.Redis(platform)
	defer common.RedisRestore(platform, rd)
	bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	cacheKey := CacheAdminPrefix + key
	_, _ = rd.Del(cacheKey).Result()
	_, _ = rd.Set(cacheKey, string(bytes), 0).Result()
	return nil
}

// Initialize 初台化缓存相关信息
func Initialize() {

	t1 := tools.TimeDebugBegin("--- 开始初始化后台缓存数据 ---")
	for _, platform := range consts.PlatformCodes {

		Banks.Load(platform)          // 银行信息
		DepositCards.Load(platform)   // 存款卡
		HelpCategories.Load(platform) // 帮助分类
		UserLevels.Load(platform)     // 用户等级
		Configs.Load(platform)        // 配置
		GameVenues.Load(platform)     // 游戏场馆

		Menus.Load(platform)             // *** 菜单加载要在角色加载之前 ***
		AdminRoles.Load(platform)        // *** 要在菜单之后, 顺序不可变更 ***
		PaymentThirds.Load(platform)     // 支付平台
		UserTagCategories.Load(platform) // 用户标签分类
		PermissionIps.Load(platform)     // 授权IP
		RiskConditions.Load(platform)    // 风控条件

		// 包网平台设置
		Platforms.Load()           // 平台
		PlatformSites.Load()       // 平台 - 盘口/站点
		PlatformSiteConfigs.Load() // 平台 - 盘口/站点 - 配置

		// ip地区库加载
		IPCached.Init(platform)
	}
	_ = tools.TimeDebugAt(t1, "||| 后台缓存数据初始化结束 |||")
}
