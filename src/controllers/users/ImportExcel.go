package users

import (
	"errors"
	"fmt"
	"regexp"
	common "sports-common"
	"sports-common/config"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req"
	"github.com/xuri/excelize/v2"
)

func (ths *Users) ImportExcel(c *gin.Context) {
	name := "import_excel_file"
	file, err := c.FormFile(name)
	if err != nil {
		response.Err(c, "缺少要上传的文件")
		return
	}
	fileName := file.Filename
	extArr := strings.Split(fileName, ".")
	extLen := len(extArr)
	if extLen < 2 {
		response.Err(c, "缺少文件扩展名")
		return
	}
	extName := extArr[extLen-1]
	if extName != "xlsx" {
		response.Err(c, "文件扩展名不正确")
		return
	}
	if file.Size > 1024000 {
		response.Err(c, "文件过大")
		return
	}

	saveFileName := tools.RandString(16) + "." + extName
	oldFileName := extArr[0]
	isLetterAndNumber, _ := regexp.MatchString(`^[A-Za-z0-9]+$`, oldFileName)
	if isLetterAndNumber { //只有字母和数字 640x960,1242x2208,1440x2560
		saveFileName = tools.RandString(16) + "-" + oldFileName + "." + extName
	}
	saveFilePath := "/tmp/_import_" + saveFileName
	err = c.SaveUploadedFile(file, saveFilePath)
	if err != nil {
		response.Err(c, "保存上传文件失败")
		return
	}

	f, err := excelize.OpenFile(saveFilePath)
	if err != nil || f == nil { // 如果读取有错
		response.Err(c, "读取文件出错")
		return
	}

	// 整理总代,代理
	topAgents := map[string]string{}
	agents := map[string]string{}
	rows, _ := f.GetRows("Sheet1")
	tArr := map[string]string{}
	var uArr []UserInfoImport
	for rk, row := range rows {
		if rk <= 1 { // 跳过前2行
			continue
		}
		rowCount := len(row)
		weChat := ""
		qq := ""
		mail := ""
		phone := ""
		if rowCount >= 20 {
			weChat = row[19]
		}
		if rowCount >= 19 {
			qq = row[18]
		}
		if rowCount >= 18 {
			mail = row[17]
		}
		if rowCount >= 17 {
			phone = row[16]
		}

		u := UserInfoImport{
			UserName:  row[1],
			State:     1,
			RealName:  row[3], // 真实姓名
			Phone:     phone,
			Mail:      mail,
			QQ:        qq,
			WeChat:    weChat,
			AgentTop:  row[8], // 顶级代理
			Agent:     row[9], // 代理
			Created:   tools.GetMicroTimeStampByString(row[10]),
			LastLogin: tools.GetTimeStampByString(row[11]),
		}
		if row[1] == "正常" {
			u.State = 2
		}
		if row[6] != "" && row[6] != "0" { // 余额
			if val, err := strconv.ParseFloat(row[6], 64); err == nil {
				u.Account = val
			}
		}

		tArr[u.Agent] = u.AgentTop
		topAgents[u.AgentTop] = u.UserName
		agents[u.Agent] = u.UserName
		uArr = append(uArr, u) // 所有原始用户信息
	}

	if len(uArr) == 0 || len(topAgents) == 0 || len(agents) == 0 {
		response.Err(c, "导入失败: 文件格式或者内容有误")
		return
	}
	fmt.Println("用户总计: ", len(uArr), ", 总代总计: ", len(topAgents), ", 代理总计: ", len(agents))

	tCount := 0 // 总代数量 - 成功
	aCount := 0 // 代理数量 - 成功
	uCount := 0 // 用户数量 - 成功
	platform := request.GetPlatform(c)
	dbSession := common.Mysql(platform)
	defer dbSession.Close()

	// -- 关于总代的处理 ------------------------------------------------------------
	// 判断总代是否存在
	getSlice := func(arr map[string]string) []string {
		var rArr []string
		for k := range arr {
			rArr = append(rArr, k)
		}
		return rArr
	}
	sql := "SELECT id, username FROM users WHERE username IN ('" + strings.Join(getSlice(topAgents), "','") + "')"
	rs, err := dbSession.QueryString(sql)
	if err != nil {
		response.Err(c, "获取总代信息出错")
		return
	}
	for _, r := range rs {
		delete(topAgents, r["username"])
	}
	fTopAgents := getSlice(topAgents)

	// 添加总代
	req.SetTimeout(time.Second * 10)
	header := req.Header{"Accept": "application/json"}
	registerAgentURL := config.Get("internal.internal_admin_service") + "/agents/insert?platform=" + platform
	commissionPlan := func() string {
		rows, err := dbSession.QueryString("SELECT agent_commission FROM agent_commission_plans ORDER BY id DESC LIMIT 1")
		if err == nil && len(rows) > 0 {
			return rows[0]["agent_commission"]
		}
		return "默认"
	}()
	registerAgent := func(userName string) error {
		password := tools.RandString(8) // 生成随机密码
		param := req.Param{
			"username":         userName,
			"agent_commission": commissionPlan,
			"agent_type":       "0",
			"id":               "0",
			"label":            "转代",
			"password":         password,
			"re_password":      password,
		}
		body, err := req.Post(registerAgentURL, header, param)
		if err != nil {
			fmt.Println("系统异常: 内部调用注册代理失败:", err)
			return errors.New("系统异常: 内部调用注册代理失败")
		}

		result := response.RespInfo{} // 内部调用返回信息
		err = body.ToJSON(&result)
		if err != nil {
			fmt.Println("系统错误: 内部调用注册代理失败|", err)
			return errors.New("系统错误: 内部调用注册代理失败|")
		} else if result.ErrCode != 0 {
			return errors.New("系统错误: " + result.Message)
		}
		return nil
	}

	// 添加代理 - 顶级代理
	for _, userName := range fTopAgents {
		if err := registerAgent(userName); err != nil {
			fmt.Println("注册顶级代理出错:", err, userName)
			response.Err(c, "系统错误: 内部调用注册代理错误/"+err.Error()+"("+userName+")")
			return
		}
		tCount += 1
	}

	// --------- 关于代理的处理 -----------------------------------
	// 判断代理是否存在
	sql = "SELECT id, username FROM users WHERE username IN ('" + strings.Join(getSlice(agents), "','") + "')"
	rs, err = dbSession.QueryString(sql)
	if err != nil {
		response.Err(c, "获取总代信息出错")
		return
	}
	for _, r := range rs {
		delete(agents, r["username"])
	}
	fAgents := getSlice(agents)

	// 添加代理 - 二级代理
	iArr := map[string]int{}
	getTopId := func(topName string) int {
		if val, exists := iArr[topName]; exists {
			return val
		}
		sql := fmt.Sprintf("SELECT id FROM users WHERE username = '%s' LIMIT 1", topName)
		rows, err := dbSession.QueryString(sql)
		if err == nil && len(rows) > 0 {
			id, _ := strconv.Atoi(rows[0]["id"])
			iArr[topName] = id
			return id
		}
		return 1000
	}
	// 变更代理信息
	for _, userName := range fAgents {
		if err := registerAgent(userName); err != nil {
			fmt.Println("注册顶级代理出错:", err, userName)
			response.Err(c, "系统错误: 内部调用注册代理失败/"+err.Error()+"("+userName+")")
			return
		}
		if topName, exists := tArr[userName]; exists {
			topId := getTopId(topName)
			sql := fmt.Sprintf("UPDATE users SET top_id = %d, top_name = '%s' WHERE username = '%s'", topId, topName, userName)
			_, _ = dbSession.Exec(sql)
			aCount += 1
		} else {
			fmt.Println("严重错误: 查找代理", userName, "的上级代理失败")
		}
	}

	// -- 关于会员的处理 ----------------------------------------------------
	// 添加会员 - 注册会员
	registerUserURL := config.Get("admin_register.service") + config.Get("internal_api.register_url") // 会员注册地址
	registerUser := func(u UserInfoImport) error {
		password := tools.MD5(tools.RandString(8))
		topId := getTopId(u.Agent)
		// param := req.Param{
		// 	"user_name":    u.UserName,
		// 	"password":     password,
		// 	"confirm_pwd":  password,
		// 	"a_code":       topId,
		// 	"platform":     platform,
		// 	"device_id":    "system",
		// 	"register_id":  1,
		// 	"register_url": "https://" + c.Request.Host,
		// }
		param := req.Param{
			"user_name":         u.UserName,
			"password":          password,
			"confirm_pwd":       password,
			"a_code":            topId,
			"device_id":         "pc",
			"register_id":       1,
			"register_url":      "https://" + c.Request.Host,
			"platform":          platform,
			"withdraw_password": tools.MD5(tools.RandString(8)),
			"phone":             u.Phone,
			"realname":          u.RealName,
			"qq":                u.QQ,
			"we_chat":           u.WeChat,
		}
		body, err := req.Post(registerUserURL+"?platform="+platform, header, param)
		if err != nil {
			return errors.New("系统异常: 内部调用注册用户失败/" + err.Error() + "!(" + u.UserName + ")")
		}

		result := response.RespInfo{} // 内部调用返回信息
		err = body.ToJSON(&result)
		if err != nil {
			return errors.New("系统错误: 内部调用注册用户失败~(" + u.UserName + ")")
		}
		if result.ErrCode != 0 {
			return errors.New("系统错误: 内部调用注册用户失败/" + result.Message + "(" + u.UserName + ")")
		}
		return nil
	}

	for _, u := range uArr {
		if err := registerUser(u); err != nil {
			response.Err(c, err.Error())
			return
		}

		_ = dbSession.Begin()
		sql := "UPDATE users SET "
		data := map[string]interface{}{
			"birthday":      u.Birth,
			"email":         u.Mail,
			"created":       u.Created,
			"last_login_at": u.LastLogin,
			// "remark":        u.Remark,
		}
		for k, v := range data {
			sql += fmt.Sprintf("%s = '%v', ", k, v)
		}

		sql += fmt.Sprintf("remark = '%s' WHERE username = '%s' LIMIT 1", u.Remark, u.UserName)
		if _, err := dbSession.Exec(sql); err != nil {
			_ = dbSession.Rollback()
			fmt.Println("写入用户信息出错:", err)
		}
		sql = fmt.Sprintf("UPDATE accounts SET balance = %.2f, available = %.2f WHERE username = '%s' LIMIT 1", u.Account, u.Account, u.UserName)
		if _, err := dbSession.Exec(sql); err != nil {
			_ = dbSession.Rollback()
			fmt.Println("写入用户账户信息出错:", err)
		}

		uCount += 1
		_ = dbSession.Commit()
	}

	fmt.Println("成功注册, 总代:", tCount, ", 代理:", aCount, ", 用户:", uCount)
	response.Ok(c)
}
