package filters

import (
	"sports-admin/caches"
	"strconv"
	"strings"

	"github.com/flosch/pongo2"
)

// 会员等级
func userVipLevel(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	platform, ID := getPlatformValue(in)
	if platform == "" {
		return pongo2.AsValue("-"), nil
	}

	getUserLevel := func(value int32) string {
		levels := caches.UserLevels.All(platform)
		userDigit := uint8(value)
		for _, v := range levels {
			if v.Digit == (userDigit - 1) {
				return v.Name
			}
		}
		return "-"
	}

	levelID := int32(ID)
	return pongo2.AsValue(getUserLevel(levelID)), nil
}

// GetUserLabels 得到会员标签
func GetUserLabels(platform, value string) string {
	lines := strings.Split(value, ";") // 拆份各个标签分类
	results := []string{}
	for _, line := range lines {
		if strings.HasPrefix(line, "warning:") {
			line = "<div class='layui-badge' style='padding: 0px; display: block;'>" + strings.Split(line, "warning:")[1] + "</div>"
			results = append(results, line)
		} else if strings.Contains(line, "|") { // 如果是 分类|标签ID 模式
			categoryTags := strings.Split(line, "|")
			if len(categoryTags) == 0 {
				continue
			}
			categoryId, err := strconv.Atoi(categoryTags[0])
			if err != nil || categoryId <= 0 {
				continue
			}
			tagIds := strings.Split(categoryTags[1], ",")
			category := caches.UserTagCategories.Get(platform, categoryId) // 获得此标签分类信息
			if category == nil {                                           // 如果无结果, 则跳过
				continue
			}
			categoryText := category.Name
			for _, tagIdStr := range tagIds {
				if tagId, err := strconv.Atoi(tagIdStr); err == nil && tagId > 0 {
					for _, tag := range category.Tags {
						if uint32(tagId) == tag.Id {
							tagText := categoryText + ":" + tag.Name
							line = "<div class='layui-badge' style='padding: 0px;'>" + tagText + "</div>"
							results = append(results, line)
						}
					}
				}
			}
		}
	}
	return strings.Join(results, " ")
}

// 会员标签
func userLabel(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	value := in.Interface().(string)
	arr := strings.Split(value, ":")
	if len(arr) < 2 {
		return pongo2.AsValue(""), nil
	}
	platform := arr[0]
	labels := strings.Join(arr[1:], "")
	return pongo2.AsValue(GetUserLabels(platform, labels)), nil
}

// 会员标签分类
// 参数格式: 平台标识:id -> tianji:1008
func userTagCategory(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
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

	tagCategoryID := uint32(ID)
	levels := caches.UserTagCategories.All(platform)
	for _, v := range levels {
		if v.Id == tagCategoryID {
			return pongo2.AsValue(v.Name), nil
		}
	}
	return pongo2.AsValue("-"), nil
}

// 会员标签
func userTag(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
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

	tagID := uint32(ID)
	levels := caches.UserTagCategories.All(platform)
	for _, v := range levels {
		if v.Id == tagID {
			return pongo2.AsValue(v.Name), nil
		}
	}
	return pongo2.AsValue("-"), nil
}
