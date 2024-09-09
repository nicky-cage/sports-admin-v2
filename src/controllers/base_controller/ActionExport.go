package base_controller

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sports-admin/caches"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	models "sports-models"
	"strconv"
	"strings"
	"time"

	//"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// FieldToDateTime 转换为时间
var FieldToDateTime = func(val string) string {
	if val == "" || val == "0" {
		return ""
	}
	timestamp, err := strconv.Atoi(val)
	if err != nil {
		fmt.Println("格式化日期时间戳错误: val =", val, ", err: ", err)
	}

	// 如果是微秒
	if timestamp > 1000000000*1000000 {
		return time.UnixMicro(int64(timestamp)).Format("2006-01-02 15:04:05")
	}

	return time.Unix(int64(timestamp), 0).Format("2006-01-02 15:04:05")
}

// FieldToStatus 转换为状态
var FieldToStatus = func(val string) string {
	if val == "1" {
		return "停用"
	} else if val == "2" {
		return "正常"
	} else {
		return "未知"
	}
}

// FieldToYesNo 转换为是否
var FieldToYesNo = func(val string) string {
	if val == "1" {
		return "是"
	}
	return "否"
}

// FieldToGender 转换为性别
var FieldToGender = func(val string) string {
	if val == "1" {
		return "男"
	} else if val == "2" {
		return "女"
	} else {
		return "未知"
	}
}

// FieldToUserVip 转换用户vip等级 -- args: 累加要减去的等级
var FieldToUserVip = func(c *gin.Context, userVip interface{}, args ...int) string {
	val, err := strconv.Atoi(fmt.Sprintf("%v", userVip))
	if err != nil {
		return ""
	}
	realVal := val - func() int {
		if len(args) >= 1 {
			return args[0]
		}
		return 1
	}()
	platform := request.GetPlatform(c)
	userLevel := caches.UserLevels.Get(platform, realVal)
	if userLevel != nil {
		return userLevel.Name
	}
	return ""
}

// FieldToPaymentType 转换支付渠道
var FieldToPaymentType = func(c *gin.Context, paymentCode interface{}) string {
	code := fmt.Sprintf("%v", paymentCode)
	platform := request.GetPlatform(c)
	return caches.PaymentThirds.Get(platform, code)
}

// ExportRawDataReset 二级导出时的原始数据重置
var ExportRawDataReset = func(rArr *[]map[string]interface{}) {
	realArr := make([]map[string]interface{}, 0)
	for _, r := range *rArr {
		tArr := make(map[string]interface{})
		for rk, rv := range r { //  忽略掉二级单独的值 - 保处理数组map部分
			if row, ok := rv.(map[string]interface{}); ok { // 有可能存在二级数组的情况
				for sk, sv := range row {
					key := rk + "." + sk
					tArr[key] = sv
				}
			} else {
				tArr[rk] = rv
			}
		}
		if len(tArr) > 0 { // 只有有内容的情况下才会合并数姐
			realArr = append(realArr, tArr)
		}
	}
	*rArr = realArr
}

// HeaderNames 头部信息
var HeaderNames = func() []string {
	fArr := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	tArr := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	rArr := []string{}
	rArr = append(rArr, fArr...)
	for _, v := range fArr {
		for _, v2 := range tArr {
			rArr = append(rArr, v+v2)
		}
	}
	return rArr
}()

// ExportHeader 字段头部映射信息
type ExportHeader struct {
	Name  string `json:"name"`
	Field string `json:"field"`
}

// ExportCache 用于保存到缓存当中的数据
type ExportCache struct {
	Fields []string `json:"fields"`
	URL    string   `json:"url"`
}

// ActionExport ActionExportExcel 导出功能
type ActionExport struct {
	Columns        []ExportHeader                                          `json:"-"` // 头部信息 - 必须
	Model          common.IModel                                           `json:"-"` // 可选 - 优先
	GetSQL         func(string, *gin.Context) string                       `json:"-"` // 如果有此条, 优先选择SQL - 可选 - 次级优先
	GetData        func([]string, *gin.Context) []map[string]interface{}   `json:"-"` // 导出的数据 - 可选 - 三级优先
	ProcessRow     func(*map[string]interface{}, *gin.Context)             `json:"-"` // 对于部分字段的回调处理 - 可选
	ProcessRawData func([]string, *[]map[string]interface{}, *gin.Context) `json:"-"` // 处理原始数据
}

// GetExportRawDataByURL 得到导出的数据
var GetExportRawDataByURL = func(ec ExportCache, c *gin.Context) []byte {
	requestURL := fmt.Sprintf("%s://%s%s", func() string {
		if strings.Contains(c.Request.Host, "admin.sports") {
			return "http" // 本地环境
		}
		return "https" // 线上环境
	}(), c.Request.Host, func() string { // 上传到服务器上之后，要改成https
		if strings.Contains(ec.URL, "?") {
			rArr := strings.Split(ec.URL, "?")
			return rArr[0]
		}
		return ec.URL
	}())
	realURL, err := url.Parse(requestURL)
	if err != nil {
		panic("错误生成URL: " + err.Error())
	}
	params := url.Values{}
	params.Set("export_excel", "1")
	if strings.Contains(ec.URL, "?") {
		qArr := strings.Split(strings.Split(ec.URL, "?")[1], "&")
		for _, v := range qArr {
			tArr := strings.Split(v, "=")
			if len(tArr) == 2 {
				params.Set(tArr[0], tArr[1])
			}
		}
	}
	realURL.RawQuery = params.Encode()
	request, err := http.NewRequest("GET", realURL.String(), nil) //建立一个请求
	if err != nil {
		log.Logger.Error("Error No NewRequest:", err)
		return nil
	}

	// 设置一些头部信息
	func() {
		cookieString, err := c.Cookie("sports_login")
		if err == nil {
			adminCookie := models.AdminCookies.FromJSON(cookieString)
			cString := url.QueryEscape(adminCookie.ToJSON())
			cookie := &http.Cookie{
				Name:     "sports_login",
				Value:    cString,
				HttpOnly: true,
				Path:     "/",
				Secure:   false,
				Domain:   c.Request.Host,
			}
			request.AddCookie(cookie)
			request.Header.Add("sports_login", cString)
		}
		request.Header.Add("X-Requested-With", "xmlhttprequest")
		request.Header.Add("Accept-Charset", "utf-8")
		request.Header.Add("Accept-Language", "ja,zh-CN;q=0.8,zh;q=0.6")
		request.Header.Add("Cache-Control", "no-cache")
		request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko)")
		request.Header.Add("Connection", "close")
	}()

	tr := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, //跳过证书验证
		MaxIdleConns:          2,
		MaxIdleConnsPerHost:   1,
		MaxConnsPerHost:       1,
		DisableKeepAlives:     true,
		ResponseHeaderTimeout: time.Second * 60,
		DialContext: (&net.Dialer{
			Timeout: time.Second * 30, //超时设置
		}).DialContext,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 30,
	}
	response, err := client.Do(request) //提交
	if err != nil {
		log.Logger.Error("Can not NewClient:", err)
		return nil
	}
	defer response.Body.Close()

	rBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Logger.Error(" *** Can not RealAllResponse:", err)
		return nil
	}
	return rBytes
}

// Export 自动导出excel
func (ths *ActionExport) Export(c *gin.Context) {
	platform := request.GetPlatform(c) // 平台识别号
	// 如果提交过来的是要导出的字段, 则加密并写入redis
	if request.IsAjax(c) { // 如果是ajax
		postedData := request.GetPostedData(c)
		fArr := []string{}
		for k, v := range postedData {
			if v != "" {
				fArr = append(fArr, k)
			}
		}
		rStruct := ExportCache{
			Fields: fArr,
			URL: func() string {
				if v, exists := postedData["export_url"]; exists {
					return v.(string)
				}
				return ""
			}(),
		}
		rBytes, err := json.Marshal(rStruct)
		if err != nil {
			response.Err(c, "序列化字段信息错误: "+err.Error())
			return
		}
		rValue := string(rBytes)
		sign := tools.MD5(rValue + time.Now().Format("2006-01-02 15:04:05"))
		rKey := platform + ":export:" + sign
		redisClient := common.Redis(platform)
		defer common.RedisRestore(platform, redisClient)
		_, _ = redisClient.Set(rKey, rValue, time.Minute*1).Result()
		response.Result(c, sign)
		return
	}

	// -- 如果是准备要直下载的信息
	// 则: 1. 检测有无model, 如果有则自动生成
	signMsg := c.DefaultQuery("sign", "")
	if signMsg != "" && len(signMsg) == 32 { // 如果是跳转过来的下载
		redisClient := common.Redis(platform)
		defer common.RedisRestore(platform, redisClient)
		rKey := platform + ":export:" + signMsg
		rValue, err := redisClient.Get(rKey).Result()
		if err != nil || rValue == "" {
			response.ErrorHTML(c, "错误: 未找到原始数据或者缺少权限, 无法下载文件 ...")
			return
		}

		rCache := ExportCache{}
		err = json.Unmarshal([]byte(rValue), &rCache)
		if err != nil {
			response.ErrorHTML(c, "错误: 反序列化失败或者IP权限("+rValue+")")
			return
		}

		// 拿到需要的数据
		rows, err := func() ([]map[string]interface{}, error) {
			if rCache.URL != "" { // 如果有sql先请求第三方拿到原始数据
				rawBytes := GetExportRawDataByURL(rCache, c)
				if rawBytes == nil {
					return nil, errors.New("无法获取原始数据")
				}
				resp := response.RespInfo{}
				err := json.Unmarshal(rawBytes, &resp)
				if err != nil {
					fmt.Sprintln("原始数据: ", string(rawBytes))
					return nil, errors.New("反序列化输出结果错误(" + err.Error() + ")")
				}
				rawData, ok := resp.Data.([]interface{})
				if !ok {
					return nil, errors.New("数组获取有误(" + resp.Message + ")")
				}
				rArr := make([]map[string]interface{}, 0)
				for _, r := range rawData {
					rArr = append(rArr, r.(map[string]interface{}))
				}
				if ths.ProcessRawData != nil {
					ths.ProcessRawData(rCache.Fields, &rArr, c)
				}
				return rArr, nil
			}
			sql := func() string { // 没有指定sql的情况, 自动生成sql
				if ths.Model != nil { // 1. 如果有model则自动生成 - 优先
					return fmt.Sprintf("SELECT %s FROM %s", rValue, ths.Model.GetTableName())
				}
				if ths.GetSQL != nil { // 2. 如果有生成sql则执行此sql - 次级优先
					return ths.GetSQL(rValue, c)
				}
				return ""
			}()
			if sql != "" {
				dbSession := common.Mysql(platform)
				defer dbSession.Close()
				rows, err := dbSession.QueryInterface(sql)
				return rows, err
			}
			if ths.GetData == nil { // 3. 如果没有指定sql和model, 则返回指定数据 - 三级优先
				return nil, errors.New("未设定任何获取数据方式")
			}
			return ths.GetData(strings.Split(rValue, ","), c), nil
		}()
		if err != nil {
			response.ErrorHTML(c, "错误: "+err.Error()+", 没有数据可供下载 ...")
			return
		}
		if len(rows) == 0 {
			response.ErrorHTML(c, "错误: 查询结果为空, 没有数据可供下载 ...")
			return
		}

		// ---------- 写入 excel 文件并下载 --------------------------
		xlsx := excelize.NewFile() // 新建立excel对象
		sheetName := "Sheet1"
		columnNames := HeaderNames            // 写数据头部
		rArr := rCache.Fields                 // 拆分字段
		hasField := func(field string) bool { // 是否包含某个字段
			for _, k := range rArr {
				if k == field {
					return true
				}
			}
			return false
		}
		indexCount := 0
		for _, v := range ths.Columns { // 先写头部信息
			if !hasField(v.Field) {
				continue
			}
			cellName := columnNames[indexCount] + strconv.Itoa(1)
			cellValue := v.Name
			xlsx.SetCellValue(sheetName, cellName, cellValue)
			indexCount += 1
		}

		for rk, row := range rows {
			indexCount = 0
			if ths.ProcessRow != nil { // 如果可以指定处理字估
				ths.ProcessRow(&row, c)
			}
			for _, v := range ths.Columns {
				if !hasField(v.Field) {
					continue
				}
				cellName := columnNames[indexCount] + strconv.Itoa(rk+2) // cell 名字
				cellValue := row[v.Field]
				xlsx.SetCellValue(sheetName, cellName, cellValue)
				indexCount += 1
			}
		}

		fileName := time.Now().Format("20060102_150405") + ".xlsx"
		//xlsx.SetCellValue("Sheet1", "A2", "asdas") //_ = xlsx.SaveAs("./aaa.xlsx")
		c.Header("Content-Type", "application/octet-stream")
		c.Header("Content-Disposition", "attachment; filename="+fileName)
		c.Header("Content-Transfer-Encoding", "binary")
		_ = xlsx.Write(c.Writer)
	}

	if ths.Columns == nil {
		response.Err(c, "未设置头部信息")
		return
	}

	exportURL := func() string {
		originUrl := c.DefaultQuery("URL", "")
		if originUrl != "" {
			ajaxURL, err := url.QueryUnescape(originUrl)
			if err != nil {
				return ""
			}
			return ajaxURL
		}
		return ""
	}()

	viewData := response.ViewData{
		"columns":   ths.Columns,
		"actionURL": c.Request.URL.Path,
		"exportURL": exportURL,
	}
	response.Render(c, "export.html", viewData)
}
