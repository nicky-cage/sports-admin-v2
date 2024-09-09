package users

import (
	"fmt"
	common "sports-common"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LevelUpSave 保存用户vip等级修改
func (ths *Users) LevelUpSave(c *gin.Context) {
	id, iErr := strconv.Atoi(c.DefaultQuery("id", "0"))
	if iErr != nil || id <= 0 {
		response.Err(c, "编号有误")
		return
	}
	postedData := request.GetPostedData(c)
	vip, vExists := postedData["vip"]
	if !vExists {
		response.Err(c, "缺少VIP信息")
		return
	}
	vipLevel, vErr := strconv.Atoi(vip.(string))
	if vErr != nil || vipLevel > 10 || vipLevel < 0 {
		response.Err(c, "VIP有误")
		return
	}

	platform := request.GetPlatform(c)
	myClient := common.Mysql(platform)
	defer myClient.Close()
	user := models.User{}
	if exists, uErr := myClient.SQL("SELECT * From users WHERE id = ? LIMIT 1", id).Get(&user); uErr != nil {
		response.Err(c, "查找用户出错")
		return
	} else if !exists {
		response.Err(c, "不能找到用户")
		return
	}

	uSQL := fmt.Sprintf("UPDATE users SET vip = %d WHERE id = %d LIMIT 1", vipLevel, id)
	if _, rErr := myClient.Exec(uSQL); rErr != nil {
		response.Err(c, "保存vip信息有误")
		return
	}

	response.Ok(c)
}
