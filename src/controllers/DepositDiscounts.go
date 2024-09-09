package controllers

import (
	"encoding/json"
	"fmt"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

// DepositDiscount 存款优惠
type DepositDiscount = models.DepositDiscount

// OutDepositDiscount 存款优惠
type OutDepositDiscount = models.OutDepositDiscount

// DepositDiscounts 存款优惠-暂时十种支付方式
var DepositDiscounts = struct {
	List func(*gin.Context)
	Edit func(*gin.Context)
	*ActionState
	SaveDo func(*gin.Context)
}{
	List: func(c *gin.Context) { //默认首页
		depositDiscounts := make([]models.UserDepositDiscount, 0)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		if err := engine.Table("user_deposit_discounts").OrderBy("id ASC").Find(&depositDiscounts); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "获取列表错误")
			return
		}
		out := make([]OutDepositDiscount, 10)
		for i, v := range depositDiscounts {
			offerContent := []byte(v.OfferContent)
			depositDiscountInfo := DepositDiscount{}
			_ = json.Unmarshal(offerContent, &depositDiscountInfo)
			newDepositDiscount := DepositDiscount{}
			for _, vv := range depositDiscountInfo {
				if vv.Vip == "VIP9" || vv.Vip == "VIP10" {
					continue
				}
				if vv.Ratio != "0.00" && vv.DayMaxDiscount != "0.00" && vv.Multiple != "0" {
					newDepositDiscount = append(newDepositDiscount, vv)
				}
			}
			out[i].OfferContent = newDepositDiscount
			out[i].Id = v.Id
			out[i].PaymentType = v.PaymentType
			out[i].Recommend = v.Recommend
			out[i].State = v.State
			out[i].Operator = v.Operator
			out[i].Created = v.Created
			out[i].Updated = v.Updated
		}
		viewData := pongo2.Context{
			"rows": out,
		}
		viewFile := "deposit_discounts/list.html"
		if request.IsAjax(c) {
			viewFile = "deposit_discounts/_list.html"
		}
		SetLoginAdmin(c)
		response.Render(c, viewFile, viewData)
	},
	Edit: func(c *gin.Context) {
		idStr := c.DefaultQuery("id", "")
		id, _ := strconv.Atoi(idStr)
		depositDiscounts := &models.UserDepositDiscount{}
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		if _, err := engine.Table("user_deposit_discounts").Where("id=?", id).Get(depositDiscounts); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "获取信息错误")
			return
		}
		depositDiscountInfo := DepositDiscount{}
		offerContent := []byte(depositDiscounts.OfferContent)
		_ = json.Unmarshal(offerContent, &depositDiscountInfo)
		viewData := pongo2.Context{
			"r":    depositDiscounts,
			"info": depositDiscountInfo,
		}
		response.Render(c, "deposit_discounts/edit.html", viewData)
	},
	SaveDo: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		idStr := postedData["id"].(string)
		id, _ := strconv.Atoi(idStr)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		offerContent := make([]map[string]string, 0)
		for i := 0; i <= 10; i++ {
			temp := make(map[string]string)
			val, exists := postedData["VIP"+strconv.Itoa(i)+"*ratio"]
			if !exists {
				continue
			}
			ratioF, _ := strconv.ParseFloat(val.(string), 64)
			ratioStr := fmt.Sprintf("%.2f", ratioF)
			temp["ratio"] = ratioStr
			dayMaxDiscountF, _ := strconv.ParseFloat(postedData["VIP"+strconv.Itoa(i)+"*day_max_discount"].(string), 64)
			dayMaxDiscountStr := fmt.Sprintf("%.2f", dayMaxDiscountF)
			temp["day_max_discount"] = dayMaxDiscountStr
			temp["multiple"] = postedData["VIP"+strconv.Itoa(i)+"*multiple"].(string)
			temp["vip"] = postedData["VIP"+strconv.Itoa(i)+"*VIP"].(string)
			offerContent = append(offerContent, temp)
		}
		offerContentJson, _ := json.Marshal(offerContent)
		uuMap := map[string]interface{}{
			"offer_content": string(offerContentJson),
			"recommend":     postedData["recommend"],
			"operator":      GetLoginAdmin(c).Name,
			"updated":       tools.NowMicro(),
		}
		if _, err := engine.Table("user_deposit_discounts").Where("id=?", id).Update(uuMap); err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "更新失败")
			return
		}
		response.Ok(c)
	},
	ActionState: &ActionState{
		Model: models.UserDepositDiscounts,
	},
}
