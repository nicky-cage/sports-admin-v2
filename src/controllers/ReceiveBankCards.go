package controllers

import (
	"sports-admin/caches"
	"sports-admin/validations"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	models "sports-models"
	"strconv"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

// ReceiveBankCards 收款银行卡
var ReceiveBankCards = struct {
	List func(*gin.Context)
	*ActionCreate
	*ActionUpdate
	*ActionSave
	*ActionDelete
	*ActionState
}{
	List: func(c *gin.Context) {
		cond := request.GetQueryCond(c, map[string]interface{}{
			"status":    "=",
			"bank_name": "=",
		})
		limit, offset := request.GetOffsets(c)
		depositCards := make([]models.DepositCard, 0)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		total, err := engine.Table("deposit_cards").Where(cond).OrderBy("id DESC").Limit(limit, offset).FindAndCount(&depositCards)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "获取列表错误")
			return
		}
		for i, v := range depositCards {
			vipsArr := strings.Split(v.Vips, ",")
			if len(vipsArr) == 0 || v.Vips == "" {
				continue
			}

			str := ""
			for _, vv := range vipsArr {
				vvInt, _ := strconv.Atoi(vv)
				str += "VIP" + strconv.Itoa(vvInt-1) + ","
			}
			depositCards[i].Vips = strings.TrimRight(str, ",")
		}

		viewData := pongo2.Context{
			"rows":  depositCards,
			"total": total,
			"banks": caches.Banks.All(platform),
			//"floatRate": fmt.Sprintf("%.2f", models.Parameters.GetValueByFloat(platform, "deposit_float_rate", 0.02, "取款浮动汇率")),
		}
		viewFile := "receive_bank_cards/index.html"
		if request.IsAjax(c) {
			viewFile = "receive_bank_cards/_list.html"
		}
		SetLoginAdmin(c)
		response.Render(c, viewFile, viewData)
	},
	ActionCreate: &ActionCreate{
		Model:    models.DepositCards,
		ViewFile: "receive_bank_cards/edit.html",
		ExtendData: func(c *gin.Context) pongo2.Context {
			platform := request.GetPlatform(c)
			banks := caches.Banks.All(platform)
			//for i := 0; i < len(banks); i++ {
			//	if banks[i].Name == "其他银行" {
			//		banks = append(banks[:i], banks[i+1:]...)
			//		i--
			//	}
			//}
			return pongo2.Context{
				"banks": banks,
			}
		},
	},
	ActionUpdate: &ActionUpdate{
		Model: models.DepositCards,
		Row: func() interface{} {
			return &models.DepositCard{}
		},
		ViewFile: "receive_bank_cards/edit.html",
		ExtendData: func(c *gin.Context) pongo2.Context {
			platform := request.GetPlatform(c)
			return pongo2.Context{
				"banks": caches.Banks.All(platform),
			}
		},
	},
	ActionSave: &ActionSave{
		Model: models.DepositCards,
		SaveBefore: func(c *gin.Context, m *map[string]interface{}) error {
			tempVipIds := ""
			for i := 0; i <= 10; i++ {
				v, ok := (*m)["vips["+strconv.Itoa(i)+"]"]
				if ok {
					tempVipIds += v.(string) + ","
				}
			}
			tempVipIds = strings.TrimRight(tempVipIds, ",")
			(*m)["vips"] = tempVipIds
			(*m)["byname"] = (*m)["bank_code"].(string) + "-" + (*m)["bank_realname"].(string) + "-" + (*m)["bank_card"].(string)
			return nil
		},
		Validator: validations.ReceiveBankCard,
		SaveAfter: func(c *gin.Context, data *map[string]interface{}) {
			platform := request.GetPlatform(c)
			caches.DepositCards.Load(platform)
		},
	},
	ActionDelete: &ActionDelete{
		Model: models.DepositCards,
		DeleteAfter: func(c *gin.Context, data interface{}) {
			platform := request.GetPlatform(c)
			caches.DepositCards.Load(platform)
		},
	},
	ActionState: &ActionState{
		Model: models.DepositCards,
		Field: "status",
	},
}
