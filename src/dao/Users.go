package dao

import (
	"context"
	"reflect"
	common "sports-common"
	"sports-common/es"
	"sports-common/log"
	"sports-common/request"
	models "sports-models"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"xorm.io/builder"
)

// UserDetail 用户详情信息
type UserDetail struct {
	User        models.User                // 用户信息
	Agent       models.User                // 代理信息
	Cards       []models.UserCard          // 用户绑卡信息
	Money       models.Account             // 用户钱包
	Lower       string                     // 下线会员。
	Notes       []models.UserNote          // 用户备注
	VirtualRows []models.UserVirtualWallet // 用户钱包
}

type LoginIp struct {
	LastIp string `json:"last_ip"` //登录IP
}

// Users 用户信息
var Users = struct {
	Detail func(*gin.Context, int) (*UserDetail, error)
}{
	Detail: func(c *gin.Context, userId int) (*UserDetail, error) {
		platform := request.GetPlatform(c)
		myClient := common.Mysql(platform)
		defer myClient.Close()

		// 获取用户基本信息
		cond := builder.NewCond().And(builder.Eq{"id": userId})
		user := models.User{}
		if exists, err := models.Users.Find(platform, &user, cond); err != nil || !exists {
			log.Err("获取用户信息有误:", err)
			return nil, err
		}

		//从Es 获取最近的登录ip
		esIndexName := platform + "_login_logs"
		boolQuery := elastic.NewBoolQuery()
		boolQuery.Must(elastic.NewMatchQuery("user_id", user.Id))
		esClient, _ := es.GetClientByPlatform(platform)
		eRes, err := esClient.Search(esIndexName).Query(boolQuery).Sort("created_at", false).
			From(0).Size(1).Do(context.Background())
		if err != nil {
			log.Err(err.Error())
		}
		temp := LoginIp{}
		list := make([]LoginIp, 0)
		for _, item := range eRes.Each(reflect.TypeOf(temp)) {
			list = append(list, item.(LoginIp))
		}
		if len(list) > 0 {
			user.LastLoginIp = list[0].LastIp
		}

		// 详情信息
		detail := &UserDetail{
			User: user,
		}

		// 获取用户代理信息
		agent := models.User{}
		cond = builder.NewCond().And(builder.Eq{"id": userId})
		_, _ = models.Users.Find(platform, &agent, cond)
		detail.Agent = agent

		// 获取用户绑卡信息
		cards := []models.UserCard{}
		cond = builder.NewCond().And(builder.Eq{"user_id": userId})
		_ = models.UserCards.FindAllNoCount(platform, &cards, cond)
		detail.Cards = cards

		// 获取用户钱包
		money := models.Account{}
		builder.NewCond().And(builder.Eq{"user_id": userId})
		_, _ = models.Accounts.Find(platform, &money, cond)
		detail.Money = money

		// 获取用下线会员
		res, _ := myClient.QueryString("SELECT COUNT(*) AS num FROM users WHERE top_id = %d", userId)
		if len(res) > 0 {
			detail.Lower = res[0]["num"]
		} else {
			detail.Lower = "0"
		}

		// 获取用户备注
		rowNotes := []models.UserNote{}
		_ = myClient.Table("user_notes").Where("user_id = ?", userId).OrderBy("created DESC").Find(&rowNotes)
		detail.Notes = rowNotes

		// 获取用户虚拟钱包
		rowWallets := []models.UserVirtualWallet{}
		_ = myClient.Table("user_virtual_coins").Where("user_id = ?", userId).OrderBy("created DESC").Find(&rowWallets)
		detail.VirtualRows = rowWallets

		return detail, nil
	},
}
