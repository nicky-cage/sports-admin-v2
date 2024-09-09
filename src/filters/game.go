package filters

import (
	"sports-admin/caches"
	"sports-common/consts"
	"strconv"
	"strings"

	"github.com/flosch/pongo2"
)

// 游戏场馆类型
func gameVenueType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return getValueFromTypes(in, consts.GameVenueTypes, "-")
}

// 游戏场馆
func gameVenue(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if value, isType := in.Interface().(string); isType {
		arr := strings.Split(value, ":")
		if len(arr) < 2 {
			return pongo2.AsValue(""), nil
		}

		platform := arr[0]
		code := arr[1]
		venues := caches.GameVenues.All(platform)
		for _, r := range venues {
			if r.Code == code {
				return pongo2.AsValue(r.Name), nil
			}
		}
	}
	return pongo2.AsValue("-"), nil
}

// 游戏场馆类型
func activityGameType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if value, isType := in.Interface().(string); isType {
		vArr := strings.Split(value, "-")
		if len(vArr) <= 1 {
			return pongo2.AsValue(vArr[0]), nil
		}

		if gameId, err := strconv.Atoi(vArr[1]); err == nil {
			if val, exists := consts.GameVenueTypes[uint8(gameId)]; exists {
				return pongo2.AsValue(vArr[0] + " - " + val), nil
			}
		}
	}
	return pongo2.AsValue(""), nil
}

// 支持平台类型
func gameSupportPlatformType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if value, isType := in.Interface().(string); isType {
		returnArr := []string{}
		idStrArr := strings.Split(value, ",") //拆分数组
		for _, typeStr := range idStrArr {
			if typeId, err := strconv.Atoi(typeStr); err == nil {
				if typeValue, exists := consts.GamePlatformTypes[uint8(typeId)]; exists {
					returnArr = append(returnArr, typeValue)
				}
			}
		}
		return pongo2.AsValue(strings.Join(returnArr, ",")), nil
	}
	return pongo2.AsValue(""), nil
}

// 展示类型
func gameDisplayType(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if value, isType := in.Interface().(string); isType {
		returnArr := []string{}
		idStrArr := strings.Split(value, ",") //拆分数组
		for _, typeStr := range idStrArr {
			if typeId, err := strconv.Atoi(typeStr); err == nil {
				if typeValue, exists := consts.GameDisplayTypes[uint8(typeId)]; exists {
					returnArr = append(returnArr, typeValue)
				}
			}
		}
		return pongo2.AsValue(strings.Join(returnArr, ",")), nil
	}
	return pongo2.AsValue(""), nil
}
