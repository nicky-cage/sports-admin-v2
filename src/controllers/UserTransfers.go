package controllers

import (
	common "sports-common"
	"sports-common/config"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"xorm.io/builder"
)

type TransferStatus struct {
	Errcode int    `json:"errcode"`
	Message string `json:"message"`
	//Data    struct {
	//	Money       string `json:"money"`
	//	GameBalance []struct {
	//		GameCode string  `json:"game_code"`
	//		Balance  float64 `json:"balance"`
	//	} `json:"game_balance"`
	//} `json:"data"`
}

// UserTransfers 用户账户调整记录
var UserTransfers = struct {
	List                func(*gin.Context)
	CheckTransferStatus func(*gin.Context)
}{
	List: func(c *gin.Context) { //默认首页
		cond := builder.NewCond()
		var startAt int64
		var endAt int64
		if value, exists := c.GetQuery("created"); !exists {
			//不存在直接倒叙。
			//current_time := time.Now().Unix()
			////今天的时间，   今天0，0开始 到 今天240
			//startAt = current_time - current_time%86400
			//endAt = startAt + 86400
		} else {
			areas := strings.Split(value, " - ")
			if areas[0] == areas[1] {
				//说明是获取今天的
				start, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
				startAt = start.UnixMicro()
				endAt = startAt + tools.SecondToMicro(86400)
			} else {
				startAt = tools.GetMicroTimeStampByString(areas[0])
				endAt = tools.GetMicroTimeStampByString(areas[1])
			}
			cond = cond.And(builder.Gte{"created": startAt}).And(builder.Lte{"created": endAt})
		}

		username := c.DefaultQuery("username", "")
		status := c.DefaultQuery("status", "")
		transferType := c.DefaultQuery("transfer_type", "")
		transferOutAccount := c.DefaultQuery("transfer_out_account", "")
		transferInAccount := c.DefaultQuery("transfer_in_account", "")
		if len(username) > 0 {
			cond = cond.And(builder.Eq{"username": username})
		}
		if len(status) > 0 {
			cond = cond.And(builder.Eq{"status": status})
		}
		if len(transferType) > 0 {
			cond = cond.And(builder.Eq{"transfer_type": transferType})
		}
		if len(transferOutAccount) > 0 {
			cond = cond.And(builder.Eq{"transfer_out_account": transferOutAccount})
		}
		if len(transferInAccount) > 0 {
			cond = cond.And(builder.Eq{"transfer_in_account": transferInAccount})
		}
		// 过滤金额为0的记录
		cond = cond.And(builder.Neq{"money": 0})

		limit, offset := request.GetOffsets(c)
		userTransfers := make([]models.UserTransfer, 0)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		total, err := engine.Table("user_transfers").Where(cond).OrderBy("created DESC").Limit(limit, offset).FindAndCount(&userTransfers)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "获取列表错误")
			return
		}

		viewData := pongo2.Context{
			"rows":  userTransfers,
			"total": total,
		}

		viewFile := "user_changes/user_transfers.html"
		if request.IsAjax(c) {
			viewFile = "user_changes/_user_transfers.html"
		}
		response.Render(c, viewFile, viewData)
	},
	CheckTransferStatus: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		idStr, exists := postedData["id"].(string)
		if !exists || exists && idStr == "0" {
			response.Err(c, "id为空")
			return
		}
		id, _ := strconv.Atoi(idStr)
		transferInfo := &models.UserTransfer{}
		platform := request.GetPlatform(c)
		_, _ = models.UserTransfers.FindById(platform, id, transferInfo)

		req.SetTimeout(50 * time.Second)
		req.Debug = true
		header := req.Header{
			"Accept": "application/json",
		}
		param := req.Param{
			"transfer_billno": transferInfo.BillNo,
		}

		baseTransferUrl := config.Get("internal.internal_game_service") + config.Get("internal_api.check_transfer_status_url")
		TransferUrl := baseTransferUrl + "?user_id=" + strconv.Itoa(int(transferInfo.UserId)) + "&platform=" + platform
		r, err := req.Post(TransferUrl, header, param)
		if err != nil {
			log.Logger.Error(err.Error())
			response.Err(c, "系统异常")
			return
		}
		transferResp := TransferStatus{}
		err = r.ToJSON(&transferResp)
		if err != nil {
			log.Err(err.Error())
			response.Err(c, "系统异常")
			return
		}

		//如果转账不成功。 有2中可能，。1未处理， 2失败，返回未处理。 应该是接口那边处理表的状态。
		response.Result(c, transferResp.Errcode)
	},
}
