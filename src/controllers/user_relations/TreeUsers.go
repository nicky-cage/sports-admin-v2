package user_relations

import (
	"encoding/json"
	"fmt"
	"sports-admin/controllers"
	common "sports-common"
	"sports-common/log"
	"sports-common/request"
	"sports-common/response"
	"sports-common/tools"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
)

// UserRelated 用户关系
type UserRelated struct {
	UserId      int           `json:"user_id"`                   // 用户编号
	UserName    string        `json:"user_name" xorm:"username"` // 用户名称
	IsAgent     int           `json:"is_agent"`                  // 是否代理
	TopId       int           `json:"top_id"`                    // 上级代理编号
	TopName     string        `json:"top_name"`                  // 上级代理名称
	Created     int           `json:"created"`                   // 注册时间
	LastLoginAt int           `json:"last_login_at"`             // 最后登录
	Users       []UserRelated `json:"users" xorm:"-"`            // 下级用户
}

// UserRelatedJson JSON输出格式
type UserRelatedJson struct {
	Id       int               `json:"id"`       // 编号
	Title    string            `json:"title"`    // 标题
	Field    string            `json:"field"`    // 字段名称, 必须唯一
	ShowLine bool              `json:"showLine"` // 是否显示连接线
	Href     string            `json:"href"`     // 跳转地址
	IsJump   bool              `json:"isJump"`   // 是否开启跳转
	Children []UserRelatedJson `json:"children"` // 下级
}

// TreeUsers 树形用户关系
func (ths *userRelations) TreeUsers(c *gin.Context) {
	platform := request.GetPlatform(c)

	var userID = 0
	if v, err := strconv.Atoi(c.DefaultQuery("id", "0")); err != nil || v == 0 {
		response.Err(c, "")
		return
	} else {
		userID = v
	}

	myClient := common.Mysql(platform)
	defer myClient.Close()

	totalUser := 0
	totalIn := 0.0
	totalOut := 0.0

	rows := ths.getTreeUsers(myClient, userID, &totalUser, &totalIn, &totalOut)
	JSONUsers, err := json.Marshal(rows)
	if err != nil {
		response.Err(c, "格式化结果出错")
		return
	}

	response.Render(c, "users/down_users.html", controllers.ViewData{
		"json_users":     string(JSONUsers),
		"total_user":     totalUser,
		"total_deposit":  totalIn,
		"total_withdraw": totalOut,
	})
}

// GetTreeUsers 得到下级用户
func (ths *userRelations) getTreeUsers(myClient *xorm.Session, userID int, totalUser *int, totalIn, totalOut *float64) []UserRelatedJson {

	list := []UserRelatedJson{} // 关系信息
	rows := []UserRelated{}     // 用户信息

	sql := fmt.Sprintf("SELECT id AS user_id, username, top_id, top_name, last_login_at, created "+
		"FROM users WHERE top_id = %d", userID)
	if err := myClient.SQL(sql).Find(&rows); err != nil {
		log.Logger.Error("获取下级用户信息失败: ", err)
	}
	if len(rows) == 0 {
		return list
	}

	iArr := []string{}
	for _, r := range rows {
		iArr = append(iArr, strconv.Itoa(r.UserId))
	}
	ids := strings.Join(iArr, ",")
	sqlTotal := strings.Join([]string{
		fmt.Sprintf("(SELECT 1 AS type, user_id, SUM(money) AS total FROM user_deposits WHERE status = 2 AND user_id IN (%s) GROUP BY user_id)", ids),
		fmt.Sprintf("(SELECT 2 AS type, user_id, SUM(money) As total FROM user_withdraws WHERE status = 2 AND user_id IN (%s) GROUP BY user_id)", ids),
	}, " UNION ALL ")
	totalRows := []struct {
		Type   int     `json:"type"`
		Total  float64 `json:"total"`
		UserId int     `json:"user_id"`
	}{}
	if err := myClient.SQL(sqlTotal).Find(&totalRows); err != nil {
		log.Logger.Error("获取用户统计信息出错:", err)
	}

	*totalUser += len(rows)
	for _, r := range rows {
		downUsers := []UserRelatedJson{}
		depositTotal := 0.0
		withdrawTotal := 0.0

		if r.IsAgent == 2 { // 如果是代理
			downUsers = ths.getTreeUsers(myClient, r.UserId, totalUser, totalIn, totalOut)
		}

		for _, t := range totalRows {
			if t.UserId != r.UserId {
				continue
			}
			if t.Type == 1 {
				depositTotal = t.Total
			} else if t.Type == 2 {
				withdrawTotal = t.Total
			}
		}

		registerTime := tools.Unix(int64(r.Created)).Format("2006-01-02 15:04:05")
		lastLoginAt := tools.Unix(int64(r.LastLoginAt)).Format("2006-01-02 15:04:05")
		title := fmt.Sprintf("<pre>%-16s - %d (充值:%-10.2f | 提款:%-10.2f | 注册时间:%s | 最后登录:%s)<pre>",
			r.UserName, r.UserId, depositTotal, withdrawTotal, registerTime, lastLoginAt)
		lRow := UserRelatedJson{
			Id:       r.UserId,
			Title:    title,
			Field:    fmt.Sprintf("user_%d_%d", r.TopId, r.UserId),
			IsJump:   true,
			Href:     "/users/detail?id=",
			Children: downUsers,
		}

		*totalIn += depositTotal
		*totalOut += withdrawTotal
		list = append(list, lRow)
	}

	return list
}
