package controllers

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"sports-admin/validations"
	common "sports-common"
	"sports-common/config"
	"sports-common/consts"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"xorm.io/builder"
)

// DividendManagements 红利管理
var DividendManagements = struct {
	Index        func(*gin.Context)
	SubmitDo     func(*gin.Context)
	FileDownload func(*gin.Context)
}{
	Index: func(c *gin.Context) { //默认首页
		downExcelFile := config.Get("internal.img_host_backend", "") + "/uploads/Excel/dividend.xlsx"
		uploadExcelFile := config.Get("internal.img_host_backend", "")
		gameVenues := make([]models.GameVenue, 0)
		platform := request.GetPlatform(c)
		engine := common.Mysql(platform)
		defer engine.Close()
		if err := engine.Table("game_venues").Where("pid=? and maintain=? and code!=?", 0, 1, "CENTERWALLET").Find(&gameVenues); err != nil {
			log.Logger.Error(err.Error())
		}
		SetLoginAdmin(c)
		response.Render(c, "dividend_managements/index.html", ViewData{
			"down_excel_url":    downExcelFile,
			"upload_excel_file": uploadExcelFile,
			"game_venus":        gameVenues,
		})
	},
	FileDownload: func(c *gin.Context) {
		filename := "dividend.xlsx"
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Writer.Header().Add("Content-Type", "application/octet-stream")
		c.File("./uploads/Excel/" + filename)
	},
	SubmitDo: func(c *gin.Context) {
		postedData := request.GetPostedData(c)
		platform := request.GetPlatform(c)
		dbSession := common.Mysql(platform)
		defer dbSession.Close()
		dividend := &models.UserDividend{}
		if postedData["operation_type"].(string) == "1" { //批量发放
			//读取上传的EXECL
			uploadPath := postedData["upload_excel"].(string)
			if len(uploadPath) <= 0 {
				response.Err(c, "请上传要导入的Excel")
				return
			}

			url := config.Get("internal.img_host_backend", "") + uploadPath
			resp, err := http.Get(url)
			if err != nil {
				log.Logger.Error(err.Error())
				response.Err(c, "系统异常")
				return
			}
			defer resp.Body.Close()

			fileDir := ""
			if consts.LogPath == "./" {
				fileDir = consts.LogPath + "tmp/" + fileDir
			} else {
				fileDir = consts.LogPath + "/tmp/" + fileDir
			}
			_, serr := os.Stat(fileDir)
			if serr != nil {
				merr := os.MkdirAll(fileDir, os.ModePerm)
				if merr != nil {
					log.Logger.Error(merr.Error())
					response.Err(c, "系统异常")
					return
				}
			}
			tempFileStr := strconv.Itoa(int(time.Now().Unix()))
			tempFilePath := tempFileStr + ".xlsx"
			out, _ := os.OpenFile(fileDir+tempFilePath, os.O_RDWR|os.O_CREATE, 0766)
			wt := bufio.NewWriter(out)
			defer out.Close()
			n, err := io.Copy(wt, resp.Body)
			if err != nil || n <= 0 {
				log.Logger.Error(err.Error())
				return
			}
			_ = wt.Flush()

			xlFile, err := xlsx.OpenFile(fileDir + tempFilePath)
			if err != nil {
				log.Logger.Error(err.Error())
				response.Err(c, "打开Excel文件错误")
				return
			}
			platform := request.GetPlatform(c)
			engine := common.Mysql(platform)
			defer engine.Close()
			//先检测--用户是否存在,是否有同名的多条记录
			needCheckRecords := make([]string, 0)
			for _, sheet := range xlFile.Sheets {
				for i, row := range sheet.Rows {
					if i == 0 {
						continue
					}
					if len(row.Cells) != 2 {
						response.Err(c, "上传的Excel格式有误")
						return
					}
					username := strings.TrimSpace(row.Cells[0].String())
					money, _ := strconv.Atoi(row.Cells[1].String())
					if money <= 0 {
						response.Err(c, "部分用户信息与金额有误，请检查!")
						return
					}
					user := &models.User{}
					has, err := engine.Table("users").Where("username=?", username).Get(user)
					if err != nil {
						log.Logger.Error(err.Error())
						response.Err(c, "系统异常")
						return
					}
					if !has {
						response.Err(c, "部分用户信息与金额有误，请检查!")
						return
					}
					if len(needCheckRecords) > 0 {
						for _, v := range needCheckRecords {
							if username == v {
								response.Err(c, "同一用户有多条记录，请检查!")
								return
							}
						}
					}
					needCheckRecords = append(needCheckRecords, username)
				}
			}
			//再插入数据
			session := dbSession
			sum := 0
			// 遍历sheet页读取
			for _, sheet := range xlFile.Sheets {
				//遍历行读取
				for i, row := range sheet.Rows {
					if i == 0 {
						continue
					}
					if len(row.Cells) != 2 {
						response.Err(c, "上传的Excel格式有误")
						return
					}
					username := strings.TrimSpace(row.Cells[0].String())
					money := strings.TrimSpace(row.Cells[1].String())
					user := &models.User{}
					_, uerr := session.Table("users").Where("username=?", username).Get(user)
					if uerr != nil {
						log.Logger.Error(uerr.Error())
						_ = session.Rollback()
						response.Err(c, "系统异常")
						return
					}
					flowMultiple, _ := strconv.Atoi(postedData["flow_multiple"].(string))
					iMap := map[string]interface{}{
						"bill_no":          tools.GetBillNo("hl", 5),
						"username":         username,
						"user_id":          user.Id,
						"top_name":         user.TopName,
						"top_id":           user.TopId,
						"type":             postedData["type"].(string), //红利类型
						"venue":            postedData["venue"].(string),
						"money":            money,
						"money_type":       postedData["money_type"].(string),
						"operation_type":   postedData["operation_type"].(string),
						"flow_limit":       postedData["flow_limit"].(string),
						"flow_multiple":    flowMultiple,
						"applicant":        GetLoginAdmin(c).Name,
						"applicant_remark": postedData["applicant_remark"].(string),
						"vip":              user.Vip,
						"created":          tools.NowMicro(),
					}
					iMap["turnover_amount"] = 0.00
					if iMap["flow_limit"] == "2" {
						fMoney, _ := strconv.ParseFloat(money, 64)
						iMap["turnover_amount"] = fMoney * float64(flowMultiple)
					}
					if _, err := session.Table("user_dividends").Insert(iMap); err != nil {
						log.Logger.Error(err.Error())
						_ = session.Rollback()
						response.Err(c, "系统异常")
						return
					}
					sum++
				}
			}
			_ = session.Commit()
			response.Message(c, "申请红利，成功"+strconv.Itoa(sum)+"条，失败0条")
		} else { //单会员
			cond := builder.NewCond().And(builder.Eq{"username": postedData["usernames"]})
			user := models.User{}
			exists, err := models.Users.Find(platform, &user, cond)
			if err != nil {
				log.Logger.Error(err.Error())
				response.Err(c, "查找会员失败")
				return
			}
			if !exists {
				response.Err(c, "不存在该会员")
				return
			}
			if err := validations.CheckDividendMoney(postedData); err != nil {
				response.Err(c, err.Error())
				return
			}
			if postedData["flow_limit"].(string) == "2" {
				if err := validations.CheckDividendFlowMultiple(postedData); err != nil {
					response.Err(c, err.Error())
					return
				}
			}
			dType, _ := strconv.Atoi(postedData["type"].(string))
			dMoney, _ := strconv.ParseFloat(postedData["money"].(string), 64)
			dMoneyType, _ := strconv.Atoi(postedData["money_type"].(string))
			dOperationType, _ := strconv.Atoi(postedData["operation_type"].(string))
			dFlowLimit, _ := strconv.Atoi(postedData["flow_limit"].(string))
			dFlowMultiple, _ := strconv.Atoi(postedData["flow_multiple"].(string))
			dividend.BillNo = tools.GetBillNo("hl", 5)

			dividend.Username = postedData["usernames"].(string)
			dividend.UserId = user.Id
			dividend.TopName = user.TopName
			dividend.TopId = user.TopId

			dividend.Venue = postedData["venue"].(string)
			dividend.Type = uint8(dType)
			dividend.Money = dMoney
			dividend.MoneyType = uint8(dMoneyType)
			dividend.OperationType = uint8(dOperationType)
			dividend.FlowLimit = uint8(dFlowLimit)
			dividend.FlowMultiple = uint8(dFlowMultiple)
			dividend.Applicant = GetLoginAdmin(c).Name
			dividend.ApplicantRemark = postedData["applicant_remark"].(string)
			dividend.State = 1
			dividend.Created = tools.NowMicro()
			dividend.IsAutomatic = 1
			dividend.TurnoverAmount = 0.00
			if dFlowLimit == 2 { //有流水
				dividend.TurnoverAmount = dMoney * float64(dFlowMultiple)
			}
			if _, err := dbSession.NoAutoTime().Insert(dividend); err != nil {
				log.Logger.Error(err.Error())
				response.Err(c, "提交失败")
				return
			}
			response.Message(c, "提交成功")
		}
	},
}
