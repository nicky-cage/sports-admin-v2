package filters

import (
	"strconv"
	"strings"

	"github.com/flosch/pongo2"
)

// 过滤器映射函数
var filterFunctions = map[string]func(*pongo2.Value, *pongo2.Value) (*pongo2.Value, *pongo2.Error){
	"channel_vip":                  channelVip,
	"datetime":                     datetime,
	"datetime64":                   datetime64,
	"login_type":                   loginType,
	"bottom_type":                  bottomType,
	"bottom_content_type":          bottomContentType,
	"state_text":                   stateText,
	"platform_type":                platformType,
	"game_venue_type":              gameVenueType,
	"agent_status":                 agentStatus,
	"agent_withdraws_process_step": agentWithdrawProcessStep,
	"agent_withdraws_status":       agentWithdrawStatus,
	"game_state":                   gameState,
	"game_maintain_state":          gameMaintainState,
	"help_content_type":            helpContentType,
	"game_support_type":            gameSupportPlatformType,
	"game_display_type":            gameDisplayType,
	"checkbox_first":               checkboxFirst,
	"checkbox_second":              checkboxSecond,
	"checkbox_three":               checkboxThree,
	"device_type":                  deviceType,
	"app_platform_type":            appPlatformType,
	"menu_level_type":              menuLevelType,
	"time_type_chang":              timeTypeChang,
	"dividend_type":                dividendType,
	"win_lose_count":               winLoseCount,
	"multi_parameter":              multiParameter,
	"user_level":                   userVipLevel,
	"user_labels":                  userLabel,
	"payment_type":                 paymentType,
	"channel_type":                 channelType,
	"string_chang_int":             stringChangInt,
	"user_tag_category":            userTagCategory,
	"user_tag":                     userTag,
	"trans_type":                   transType,
	"ip_analysis":                  ipAnalysis,
	"percentage_change":            percentageChange,
	"feedback_type":                feedbackType,
	"province":                     province,
	"city":                         city,
	"district":                     district,
	"bank":                         bank,
	"venue_lower":                  venueLower,
	"game_venues_lower":            gameVenuesLower,
	"game_venue":                   gameVenue,
	"colspan_count_app":            colspanCountApp,
	"colspan_count_pc":             colspanCountPC,
	"rowspan_count_channel":        rowspanCountChannel,
	"sport_type":                   sportType,
	"bet_status":                   betStatus,
	"activity_type":                ActivityType,
	"ip_area":                      getAreaByIp,
	"user_log_module":              userLogModule,
	"user_log_type":                userLogType,
	"activity_game_type":           activityGameType,
	"platform_name":                getPlatformName,
	"platform_site_name":           getPlatformSiteName,
	"platform_wrap":                platformWrap,
	"payment_name":                 onlinePaymentName,
	"past_time":                    pastTime,
}

// InitFilters 初始化过滤函数
func InitFilters() {
	for k, f := range filterFunctions {
		_ = pongo2.RegisterFilter(k, f)
	}
}

// 得到平台名称/数字值
func getPlatformValue(in *pongo2.Value) (string, int) {
	value := in.Interface().(string)
	arr := strings.Split(value, ":")
	if len(arr) != 2 {
		return "", 0
	}
	platform := arr[0]
	ID, err := strconv.Atoi(arr[1])
	if err != nil {
		return "", 0
	}

	return platform, ID
}
