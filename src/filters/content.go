package filters

import (
	"fmt"
	"sports-admin/caches"
	common "sports-common"
	"sports-common/consts"
	"sports-common/log"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2"
)

// 从 map[uint8]string 当中获取正确数据
func getValueFromTypes(in *pongo2.Value, data map[uint8]string, defaultValue string) (*pongo2.Value, *pongo2.Error) {
	value, isType := in.Interface().(uint8)
	if !isType {
		return pongo2.AsValue(defaultValue), nil
	}
	if val, exists := data[value]; exists {
		return pongo2.AsValue(val), nil
	}
	return pongo2.AsValue(defaultValue), nil
}

// 将数字转换为时间格式
func datetime(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	realValueIn := fmt.Sprintf("%v", in.Interface())
	realVal, err := strconv.Atoi(realValueIn)
	if err != nil || realVal == 0 {
		return pongo2.AsValue(""), nil
	}
	tm := time.Unix(int64(realVal), 0)
	timeStr := tm.Format("2006-01-02 15:04:05")
	return pongo2.AsValue(timeStr), nil
}

// 将数字转换为时间格式
func datetime64(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if value, isType := in.Interface().(int64); isType && value != 0 {
		tm := time.Unix(int64(value), 0)
		timeStr := tm.Format("2006-01-02 15:04:05")
		return pongo2.AsValue(timeStr), nil
	}
	return pongo2.AsValue(""), nil
}

// 渠道vip
func channelVip(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if value, isType := in.Interface().([]string); isType {
		tm := ""
		for _, v := range value {
			tm += "VIP" + v + ","
		}
		tm = strings.TrimRight(tm, ",")
		return pongo2.AsValue(tm), nil
	}
	return pongo2.AsValue(""), nil
}

// 登录类型
func loginType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	value, isType := in.Interface().(int8)
	if !isType {
		return pongo2.AsValue("-"), nil
	}
	if val, exists := consts.UserLoginTypes[value]; exists {
		return pongo2.AsValue(val), nil
	}
	return pongo2.AsValue("-"), nil
}

// 登录类型
func getAreaByIp(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	IpStr, isType := in.Interface().(string) // value => ip 地址
	if !isType {
		return pongo2.AsValue("-"), nil
	}
	ip := strings.TrimSpace(IpStr) // value = ip地址
	if ip == "" {
		return pongo2.AsValue(""), nil
	}

	if val, exists := caches.IPCached.List[ip]; exists {
		return pongo2.AsValue(val), nil
	}

	return pongo2.AsValue(caches.IPCached.GetArea(ip, true)), nil
	// // 如果缓存当中存在此值
	// if v, exists := ipData.List[ip]; exists {
	// 	return pongo2.AsValue(v), nil
	// }

	// saveAndReturn := func(realArea string) (*pongo2.Value, *pongo2.Error) {
	// 	ipData.List[ip] = realArea // 设置到缓存当中
	// 	return pongo2.AsValue(realArea), nil
	// }

	// areas := strings.Trim(tools.GetAreaByIp(ip), "[]")
	// areaArr := strings.Split(areas, ",")
	// area := ""
	// if len(areaArr) <= 1 {
	// 	if areaArr[0] != "" {
	// 		return saveAndReturn(areaArr[0])
	// 		// return pongo2.AsValue(areaArr[0]), nil
	// 	}
	// 	return saveAndReturn("*未知地区*")
	// 	// return pongo2.AsValue("*未知地区*"), nil
	// }
	// if areaArr[0] == "本机地址" {
	// 	return saveAndReturn("-本机地址-")
	// 	// return pongo2.AsValue("本机地址"), nil
	// }
	// area += areaArr[0]
	// if len(areaArr) <= 2 || areaArr[1] == "" {
	// 	return saveAndReturn(area)
	// 	//return pongo2.AsValue(area), nil
	// }
	// area += "-" + areaArr[1]
	// if len(areaArr) <= 3 || areaArr[2] == "" {
	// 	return saveAndReturn(area)
	// 	//return pongo2.AsValue(area), nil
	// }
	// area += "-" + areaArr[2]
	// return saveAndReturn(area)
	//return pongo2.AsValue(area), nil
}

// 底部信息类型
func bottomType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if value, isType := in.Interface().(uint8); isType {
		if val, exists := consts.BottomTypes[value]; exists {
			return pongo2.AsValue(val), nil
		}
	}
	return pongo2.AsValue(""), nil
}

// 底部信息配置 - 内容类型
func bottomContentType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return getValueFromTypes(in, consts.BottomContentTypes, "-")
}

// 维护设置
func platformType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return getValueFromTypes(in, consts.PlatformTypes, "-")
}

// 代理状态
func agentStatus(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return getValueFromTypes(in, consts.AgentStatus, "-")
}

// 代理提现步聚
func agentWithdrawProcessStep(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return getValueFromTypes(in, consts.AgentWithdrawsProcessStep, "-")
}

// 代理提现状态
func agentWithdrawStatus(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return getValueFromTypes(in, consts.AgentWithdrawsStatus, "-")
}

//站点场馆过滤
func gameState(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return getValueFromTypes(in, consts.GameState, "-")
}

// 游戏维护状态
func gameMaintainState(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return getValueFromTypes(in, consts.GameMaintainState, "-")
}

//帮助文本类型
func helpContentType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return getValueFromTypes(in, consts.HelpsContentType, "-")
}

// 复选一
func checkboxFirst(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	var str, isType = in.Interface().(string)
	if !isType {
		return pongo2.AsValue(nil), nil
	}
	if strings.Count(str, "0") > 0 {
		return pongo2.AsValue("1"), nil
	}
	return pongo2.AsValue(nil), nil
}

// 复选二
func checkboxSecond(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	var str, isType = in.Interface().(string)
	if !isType {
		return pongo2.AsValue(nil), nil
	}
	if strings.Count(str, "1") > 0 {
		return pongo2.AsValue("1"), nil
	}
	return pongo2.AsValue(nil), nil
}

// 复选三
func checkboxThree(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	var str, isType = in.Interface().(string)
	if !isType {
		return pongo2.AsValue(nil), nil
	}
	if strings.Count(str, "2") > 0 {
		return pongo2.AsValue("1"), nil
	}
	return pongo2.AsValue(nil), nil
}

// 设置类型
func deviceType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return getValueFromTypes(in, consts.DeviceType, "-")
}

// app平台类型
func appPlatformType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return getValueFromTypes(in, consts.AppPlatformType, "-")
}

//时间是字符串,转换格式
func timeTypeChang(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	ac, _ := in.Interface().(string)
	va, _ := strconv.Atoi(ac)
	if va != 0 {
		tm := time.Unix(int64(va), 0)
		timeStr := tm.Format("2006-01-02 15:04:05")
		return pongo2.AsValue(timeStr), nil
	}
	return pongo2.AsValue(""), nil

}

// 红利平台类型
func dividendType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return getValueFromTypes(in, consts.DividendType, "-")
}

//存款优惠大类
func paymentType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return getValueFromTypes(in, consts.PaymentTypes, "-")
}

//代理总输赢计算 in 是提款, param 是存款,有些还需要多个参数的。
func winLoseCount(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	first, _ := in.Interface().(string)
	second, _ := param.Interface().(string)

	f, _ := strconv.ParseFloat(first, 64)
	s, _ := strconv.ParseFloat(second, 64)

	money := f - s

	return pongo2.AsValue(money), nil
}

//多参数计算，in是之前计算的结果 （提款-存款） param 是中心钱包的钱
func multiParameter(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	first, _ := in.Interface().(float64)
	second, _ := param.Interface().(string)
	s, _ := strconv.ParseFloat(second, 64)
	money := first + s
	return pongo2.AsValue(money), nil
}

// 渠道类型
func channelType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return getValueFromTypes(in, consts.ChannalType, "-")
}

// 字符转int
func stringChangInt(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	second, _ := in.Interface().(string)
	valueInt, _ := strconv.Atoi(second)
	return pongo2.AsValue(uint8(valueInt)), nil
}

// 转账类型
func transType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	value, isType := in.Interface().(int)
	if !isType {
		return pongo2.AsValue("-"), nil
	}
	if val, exists := consts.TransType[value]; exists {
		return pongo2.AsValue(val), nil
	}
	return pongo2.AsValue("-"), nil

}

// ip统计
func ipAnalysis(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	ip, _ := in.Interface().(string)
	area := tools.GetAreaByIp(ip)
	return pongo2.AsValue(area), nil
}

// 百分号
func percentageChange(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	var val float64
	per, isType := in.Interface().(float64)
	if !isType {
		perstr, _ := in.Interface().(string)
		perf, _ := strconv.ParseFloat(perstr, 64)
		val = perf * 100
	} else {
		val = per * 100
	}

	return pongo2.AsValue(val), nil
}

// 意见反馈类型
func feedbackType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return getValueFromTypes(in, consts.FeedbackTypes, "-")
}

// 省份名称
func province(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	platform, ID := getPlatformValue(in)
	if platform == "" {
		return pongo2.AsValue("-"), nil
	}

	var val string
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	var pr []models.Province
	err := dbSession.Table("provinces").Where("id = ?", ID).Find(&pr)
	if err != nil {
		log.Err(err.Error())
	}
	if len(pr) > 0 {
		val = pr[0].Name
		return pongo2.AsValue(val), nil
	}
	return pongo2.AsValue("-"), nil
}

// 城市名称
func city(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	platform, ID := getPlatformValue(in)
	if platform == "" {
		return pongo2.AsValue("-"), nil
	}

	var val string
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	var pr []models.City
	err := dbSession.Table("cities").Where("id = ?", ID).Find(&pr)
	if err != nil {
		log.Err(err.Error())
	}
	if len(pr) > 0 {
		val = pr[0].Name
		return pongo2.AsValue(val), nil
	}
	return pongo2.AsValue("-"), nil
}

// 县区名称
func district(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	platform, ID := getPlatformValue(in)
	if platform == "" {
		return pongo2.AsValue("-"), nil
	}

	var val string
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	var pr []models.District
	err := dbSession.Table("districts").Where("id = ?", ID).Find(&pr)
	if err != nil {
		log.Err(err.Error())
	}
	if len(pr) > 0 {
		val = pr[0].Name
		return pongo2.AsValue(val), nil
	}
	return pongo2.AsValue("-"), nil
}

// 获取用户银行名称
func bank(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	value := in.Interface().(string)
	arr := strings.Split(value, ":")
	if len(arr) != 2 {
		return pongo2.AsValue("-"), nil
	}
	platform := arr[0]
	ID, err := strconv.Atoi(arr[1])
	if err != nil {
		return pongo2.AsValue("-"), nil
	}

	banks := caches.Banks.All(platform)
	bankID := uint32(ID)
	for _, v := range banks {
		if v.Id == bankID {
			return pongo2.AsValue(v.Name), nil
		}
	}
	return pongo2.AsValue("-"), nil

}

// 场錧小写
func venueLower(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	platform, ID := getPlatformValue(in)
	if platform == "" {
		return pongo2.AsValue("-"), nil
	}

	topID := uint32(ID)
	var name string
	dbSession := common.Mysql(platform)
	defer dbSession.Close()
	var pr []models.GameVenue
	err := dbSession.Table("game_venues").Where("pid = ?", topID).Find(&pr)
	if err != nil {
		log.Err(err.Error())
	}

	for _, v := range pr {
		name = name + v.Name + ","
	}
	venueName := strings.Trim(name, ",")
	return pongo2.AsValue(venueName), nil
}

// 游戏场馆
func gameVenuesLower(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	value := in.Interface().(string)
	arr := strings.Split(value, ":")
	if len(arr) < 2 {
		return pongo2.AsValue("-"), nil
	}
	platform := arr[0]
	val := strings.Join(arr[1:], "")
	sArr := strings.Split(val, ",")
	if val != "" && val != " " {
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		sql := "SELECT name FROM game_venues WHERE id IN (" + strings.Join(sArr, ",") + ")"
		res, err := dbSession.QueryString(sql)
		if err != nil {
			log.Err(err.Error())
		}
		var venues string
		for _, v := range res {
			venues = venues + v["name"] + ","
		}
		val := strings.TrimRight(venues, ",")
		return pongo2.AsValue(val), nil
	}
	return pongo2.AsValue("-"), nil
}

// 体育类型
func sportType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return getValueFromTypes(in, consts.SportTypes, "-")
}

// 投注状态
func betStatus(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	value, isType := in.Interface().(int)
	if !isType {
		return pongo2.AsValue("-"), nil
	}
	if val, exists := consts.BetStatus[value]; exists {
		return pongo2.AsValue(val), nil
	}
	return pongo2.AsValue("-"), nil
}

// ActivityType 活动类型
func ActivityType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	value, isType := in.Interface().(int32)
	if !isType {
		vals, _ := in.Interface().(string)
		values, _ := strconv.Atoi(vals)
		if val, exists := consts.ActivityTypes[int32(values)]; exists {
			return pongo2.AsValue(val), nil
		}
		return pongo2.AsValue("-"), nil
	}
	if val, exists := consts.ActivityTypes[value]; exists {
		return pongo2.AsValue(val), nil
	}
	return pongo2.AsValue("-"), nil
}

// 过去多少h/m/s
func pastTime(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	realValueIn := fmt.Sprintf("%v", in.Interface())
	realVal, err := strconv.Atoi(realValueIn)
	if err != nil || realVal <= 0 {
		return pongo2.AsValue(""), nil
	}

	past := ""
	if realVal < 60 {
		past = fmt.Sprintf("%d s", realVal)
	} else if realVal < 3600 {
		pa := realVal % 60
		if pa == 0 {
			past = fmt.Sprintf("%d m", realVal/60)
		} else {
			past = fmt.Sprintf("%d m %d s", realVal/60, pa)
		}
	} else if realVal < 86400 {
		pa := realVal % 3600
		if pa == 0 {
			past = fmt.Sprintf("%d h", realVal/3600)
		} else {
			past = fmt.Sprintf("%d h %d m", realVal/3600, pa/60)
		}
	} else {
		past = fmt.Sprintf("%d d", realVal/86400)
	}
	return pongo2.AsValue(past), nil
}
