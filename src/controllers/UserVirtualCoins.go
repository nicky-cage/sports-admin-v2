package controllers

import (
	"errors"
	"fmt"
	common "sports-common"
	"sports-common/request"
	models "sports-models"

	"github.com/gin-gonic/gin"
)

// UserVirtualCoins 会员虚拟币
var UserVirtualCoins = struct {
	*ActionList
	*ActionCreate
	*ActionUpdate
	*ActionSave
	*ActionDelete
}{
	ActionList: &ActionList{
		Model:    models.UserVirtualCoins,
		ViewFile: "user_virtual_coins/list.html",
		Rows: func() interface{} {
			return &[]models.UserVirtualCoin{}
		},
		QueryCond: map[string]interface{}{
			"user_name":      "%",
			"wallet_address": "%",
		},
	},
	ActionCreate: &ActionCreate{
		Model:    models.UserVirtualCoins,
		ViewFile: "user_virtual_coins/edit.html",
	},
	ActionUpdate: &ActionUpdate{
		Model:    models.UserVirtualCoins,
		ViewFile: "user_virtual_coins/edit.html",
		Row: func() interface{} {
			return &models.UserVirtualCoin{}
		},
	},
	ActionSave: &ActionSave{
		Model: models.UserVirtualCoins,
		CreateBefore: func(c *gin.Context, m *map[string]interface{}) error {
			userName, exists := (*m)["user_name"]
			if !exists {
				return errors.New("查询相关用户出错")
			}
			if _, exists := (*m)["wallet_address"]; !exists {
				return errors.New("缺少虚拟货币付款地址")
			}
			platform := request.GetPlatform(c)
			sql := "SELECT id FROM users WHERE username = ? LIMIT 1"
			dbSession := common.Mysql(platform)
			defer dbSession.Close()
			rows, err := dbSession.QueryString(sql, userName.(string))
			if err != nil || len(rows) == 0 {
				fmt.Println("rerr = ", err)
				return errors.New("查询相关用户出错或者找到这个用户")
			}
			(*m)["user_id"] = rows[0]["id"]

			//地址重复
			idStr := (*m)["id"].(string)
			if idStr == "0" {
				platform := request.GetPlatform(c)
				dbSession := common.Mysql(platform)
				defer dbSession.Close()
				var temp models.UserVirtualCoin
				dbSession.Table("user_virtual_coins").Where("wallet_address=?", (*m)["wallet_address"].(string)).Get(&temp)
				if temp.Id > 0 {
					return errors.New("该地址已存在")
				}
			}
			return nil
		},
		UpdateBefore: func(c *gin.Context, m *map[string]interface{}) error {
			if _, exists := (*m)["wallet_address"]; !exists {
				return errors.New("缺少虚拟货币付款地址")
			}
			return nil
		},
	},
	ActionDelete: &ActionDelete{
		Model: models.UserVirtualCoins,
	},
}
