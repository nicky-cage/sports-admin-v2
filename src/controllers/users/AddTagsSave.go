package users

import (
	"sports-common/request"
	"sports-common/response"
	models "sports-models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// AddTagsSave 保存标签
func (ths *Users) AddTagsSave(c *gin.Context) {
	postedData := request.GetPostedData(c)
	userIds, exists := postedData["user_ids"]
	if !exists {
		response.Err(c, "缺少用户相关信息")
		return
	}
	userTags, exists := postedData["user_tags"]
	if !exists {
		response.Err(c, "缺少标签相关信息")
		return
	}
	userCategoryTagIds := userTags.(string)
	if userCategoryTagIds == "" {
		response.Err(c, "提交的用户标签信息有误")
		return
	}
	categoryTags := strings.Split(userCategoryTagIds, ";")
	postedTags := map[string]string{} // 重新组织并保存前端提交的标签相关信息
	for _, ct := range categoryTags {
		arr := strings.Split(ct, "|")
		key := arr[0]
		if v, exists := postedTags[key]; exists {
			postedTags[key] = v + "," + arr[1]
		} else {
			postedTags[key] = arr[1]
		}
	}

	strIds := strings.Split(userIds.(string), ",") // 拆分用户编号
	if len(strIds) == 0 {
		response.Err(c, "提交的用户信息有误")
		return
	}
	filterIds := func(ids string) string { // 过滤掉重得的用户编号
		arr := strings.Split(ids, ",")
		tempData := map[string]bool{}
		for _, v := range arr {
			tempData[v] = true
		}
		var ret []string
		for k := range tempData {
			ret = append(ret, k)
		}
		return strings.Join(ret, ",")
	}
	platform := request.GetPlatform(c)
	for _, idStr := range strIds {
		userId, err := strconv.Atoi(idStr)
		if err != nil { // 如果转换为整型出错
			response.Err(c, "用户编号信息有误")
			return
		}
		var user models.User
		if exists, err := models.Users.FindById(platform, userId, &user); err != nil || !exists {
			response.Err(c, "查找不到用户相关信息")
			return
		}

		data := map[string]interface{}{
			"id": userId,
		}
		// 处理批量标签相关
		fromList := strings.Split(user.Label, ";")
		for _, list := range fromList {
			tempArr := strings.Split(list, "|")
			categoryId := tempArr[0]
			if v, exists := postedTags[categoryId]; exists { // 表示, 两者同有
				postedTags[categoryId] = filterIds(v + "," + tempArr[1])
			} else { // 表示, 用户现在有的标签
				postedTags[categoryId] = list
			}
		}
		var ret []string
		for k, v := range postedTags {
			ret = append(ret, k+"|"+v)
		}
		data["label"] = strings.Join(ret, ";")
		if err := models.Users.Update(platform, data); err != nil {
			response.Err(c, "处理添加用户标签时出错")
			return
		}
	}
	response.Ok(c)
}
