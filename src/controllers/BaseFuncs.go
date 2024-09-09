package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	common "sports-common"
	"sports-common/config"
	"sports-common/consts"
	"sports-common/tools"
	models "sports-models"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"xorm.io/builder"
)

// OnlineUserData 在线用户数据
type OnlineUserData struct {
	UserNames  []string
	UserTokens map[string]string
	LastTime   int64
}

// OnlineUserDatas 保存用户数据 [平台识别号]用户数据
var OnlineUserDatas = struct {
	Locker sync.Mutex
	Data   map[string]OnlineUserData
}{}

// Render 通用默认输出渲染
// func Render(c *gin.Context, viewFile string, args ...ViewData) {
// 	viewData := ViewData{}
// 	if len(args) >= 1 {
// 		viewData = args[0]
// 	}
// 	viewData["STATIC_URL"] = config.Get("internal.img_host_backend", "") // 全局静态URL
// 	response.Render(c, viewFile, viewData)
// }

// GetLoginAdmin 获取登录用户session
func GetLoginAdmin(c *gin.Context) models.LoginAdmin {
	admin, err := models.LoginAdmins.GetLoginByRequest(c)
	if admin == nil {
		fmt.Println("Fatal: ", err)
		panic("Fatal: 无法获取用户相关登录信息~")
	}
	return *admin
}

// SetLoginAdmin 设置登录信息
func SetLoginAdmin(c *gin.Context) {
	c.Set("_admin", GetLoginAdmin(c))
}

// BlockedUserNames 批量禁用用户
func BlockedUserNames(platform, userNames string) error {
	names := strings.Split(userNames, ",")
	var messages []string
	for _, userName := range names {
		user := &models.User{}
		cond := builder.NewCond().And(builder.Eq{"username": userName}).And(builder.Neq{"status": 1}) // 指定用户名称、当前状态应是正常的
		if exists, err := models.Users.Find(platform, user, cond); err != nil {                       // 如果查询出错
			messages = append(messages, err.Error())
		} else if !exists { // 如果用户不存在
			messages = append(messages, "用户不存在: "+userName)
		} else if err := user.Disable(platform); err != nil { // 从数据库修改用户状态
			messages = append(messages, err.Error())
		} else if err = models.DelUserByIdHandler(platform, int(user.Id)); err != nil { // 从缓存中删除用户
			messages = append(messages, err.Error())
		}
	}

	if len(messages) > 0 {
		return errors.New(strings.Join(messages, ","))
	}
	return nil
}

// GetTokenByUserName 依据用户名称得到用户token
func GetTokenByUserName(platform, userName string) string {
	rd := common.Redis(platform)
	defer common.RedisRestore(platform, rd)
	key := consts.CacheKeyUserNameToken + userName
	token, _ := rd.Get(key).Result()
	return token
}

// GetAllOnlineUser 得到所有在线用户
func GetAllOnlineUser(platform string) []string {
	ud, exists := OnlineUserDatas.Data[platform]
	if exists && ud.LastTime+5 > tools.Now() { // 5 秒以内, 不再获取新的数据
		return ud.UserNames
	}

	OnlineUserDatas.Locker.Lock()
	defer OnlineUserDatas.Locker.Unlock()
	if !exists {
		ud = OnlineUserData{LastTime: 0}
		OnlineUserDatas.Data = map[string]OnlineUserData{
			platform: ud,
		}
	}

	rd := common.Redis(platform)
	defer common.RedisRestore(platform, rd)
	key := consts.CacheKeyUserNameToken + "*"
	result, err := rd.Keys(key).Result()
	if err != nil {
		fmt.Println("获取所有在线用户出错: ", err)
		return nil
	}

	userNames := []string{}
	userTokens := map[string]string{}

	for _, r := range result {
		sArr := strings.Split(r, ":")
		if len(sArr) > 2 {
			userName := sArr[2]
			userNames = append(userNames, userName)
			userTokens[userName] = userName
		}
	}

	ud.LastTime = tools.Now()
	ud.UserNames = userNames
	ud.UserTokens = userTokens
	OnlineUserDatas.Data[platform] = ud

	return userNames
}

// CheckGoogleCode 检验google验证码
func CheckGoogleCode(platform string, c *gin.Context, m *map[string]interface{}) error {
	loginAdmin := GetLoginAdmin(c)
	loginAdminID := loginAdmin.Id
	admin := models.Admin{}
	if exists, err := models.Admins.FindById(platform, int(loginAdminID), &admin); !exists || err != nil {
		return errors.New("检验用户信息失败")
	}
	if config.EnvIsProduct() && admin.GoogleEnable() { // 如果已启用
		code, exists := (*m)["google_code"]
		if !exists {
			return errors.New("必须输入谷歌验证密码")
		}
		if chk, err := tools.NewGoogleAuth().VerifyCode(admin.GoogleSecret, code.(string)); err != nil {
			return err
		} else if !chk {
			return errors.New("谷歌验证密码错误")
		}
	}
	return nil
}

// HttpGet 调用外部url
func HttpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))
	return string(body), err
}
