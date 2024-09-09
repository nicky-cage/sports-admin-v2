package controllers

import (
	"errors"
	"fmt"
	common "sports-common"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ReceiveVirtuals 收款usdt
var ReceiveVirtuals = struct {
	SetFloatRate func(*gin.Context)
	*ActionList
	*ActionCreate
	*ActionUpdate
	*ActionSave
	*ActionDelete
	*ActionState
}{
	ActionList: &ActionList{
		Model:    models.DepositVirtuals,
		ViewFile: "receive_virtuals/list.html",
		Rows: func() interface{} {
			return &[]models.DepositVirtual{}
		},
		QueryCond: map[string]interface{}{
			"wallet_type":    "%",
			"wallet_address": "=",
			"state":          "=",
			"auto_rate":      "=",
		},
	},
	ActionCreate: &ActionCreate{
		Model:    models.DepositVirtuals,
		ViewFile: "receive_virtuals/edit.html",
		ExtendData: func(c *gin.Context) ViewData {
			return ViewData{
				"rate": tools.GetExchangeRate(),
			}
		},
	},
	ActionUpdate: &ActionUpdate{
		Model: models.DepositVirtuals,
		Row: func() interface{} {
			return &models.DepositVirtual{}
		},
		ViewFile: "receive_virtuals/edit.html",
		ExtendData: func(c *gin.Context) ViewData {
			return ViewData{
				"rate": tools.GetExchangeRate(),
			}
		},
	},
	ActionSave: &ActionSave{
		Model: models.DepositVirtuals,
		SaveBefore: func(c *gin.Context, m *map[string]interface{}) error {
			idStr := (*m)["id"].(string)
			if idStr == "0" {
				platform := request.GetPlatform(c)
				dbSession := common.Mysql(platform)
				defer dbSession.Close()
				var temp models.DepositVirtual
				dbSession.Table("deposit_virtuals").Where("wallet_address=?", (*m)["wallet_address"].(string)).Get(&temp)
				if temp.Id > 0 {
					return errors.New("该地址已存在")
				}
			}
			return nil
		},
		SaveAfter: func(c *gin.Context, m *map[string]interface{}) {
			stateIn, exists := (*m)["state"]
			if !exists {
				return
			}

			state, err := strconv.Atoi(stateIn.(string))
			if err != nil || (state != 1 && state != 2) {
				fmt.Println("----------- 状态获取的值不正确 --------")
				return
			}

			if state == 2 { // 如果要转为启用
				idStr := (*m)["id"].(string)
				id, err := strconv.Atoi(idStr)
				if err != nil {
					return
				}
				platform := request.GetPlatform(c)
				dbSession := common.Mysql(platform)
				defer dbSession.Close()

				if id > 0 {
					sql := fmt.Sprintf("SELECT wallet_type FROM deposit_virtuals WHERE id = %d", id)
					rows, err := dbSession.QueryString(sql)
					if err != nil || len(rows) == 0 {
						fmt.Println("获取数据有误或者不需要更新: ", err)
						return
					}
					walletType, _ := strconv.Atoi(rows[0]["wallet_type"])
					sql = fmt.Sprintf("UPDATE deposit_virtuals SET state = 1 WHERE state = 2 AND wallet_type = %d AND id <> %d", walletType, id)
					dbSession.Exec(sql)
					return
				}

				walletType, _ := strconv.Atoi((*m)["wallet_type"].(string))
				walletAddress := (*m)["wallet_address"].(string)
				dbSession.Exec("UPDATE deposit_virtuals SET state = 1 WHERE wallet_type = ? AND wallet_address <> ?", walletType, walletAddress)
			}
		},
	},
	ActionDelete: &ActionDelete{
		Model: models.DepositVirtuals,
	},
	ActionState: &ActionState{
		Model: models.DepositVirtuals,
	},
	SetFloatRate: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		fRate, exists := postedData["float_rate"]
		if !exists {
			response.Err(c, "设置浮动汇率失败: 未设置值")
			return
		}
		floatRate, err := strconv.ParseFloat(fRate.(string), 64)
		if err != nil || floatRate > 0.5 || floatRate < 0.0 {
			response.Err(c, "设置浮动汇率失败")
			return
		}

		platform := request.GetPlatform(c)
		fmt.Println("float = ", floatRate)
		err = models.Parameters.SetValue(platform, "deposit_float_rate", floatRate, 0, "存款浮动汇率")
		if err != nil {
			response.Err(c, "设置存款浮动汇率失败")
			return
		}

		response.Ok(c)
	},
}
