package dao

import (
	"fmt"
	"sports-admin/caches"
	common "sports-common"
	"sports-common/config"
	"sports-common/log"
	"strconv"
	"time"

	"github.com/imroc/req"
)

// UserGameAccount 用户游戏账户信息
type UserGameAccount struct {
	Name    string `json:"name"`    // 游戏名称1
	Code    string `json:"code"`    // 游戏代码
	Account string `json:"account"` // 账户信息
}

// GameBalance 游戏账户余额
type GameBalance struct {
	Errcode int    `json:"errcode"`
	Message string `json:"message"`
	Data    struct {
		GameCode string  `json:"game_code"`
		Balance  float64 `json:"balance"`
	}
}

// UserGames 用户 - 游戏场馆 - 资金
var UserGames = struct {
	GetAccounts func(string, int) []UserGameAccount // 得到此用户所有游戏账户
}{
	GetAccounts: func(platform string, userId int) []UserGameAccount {
		gameVenues := caches.GameVenues.All(platform) // 得到所有游戏场馆
		//cond := builder.NewCond().And(builder.Eq{"user_id": userId})
		//gRows := []models.UserGame{}
		//_ = models.UserGames.FindAllNoCount(&gRows, cond)
		//从接口中获取所有场馆的额度
		req.SetTimeout(20 * time.Second)
		req.Debug = false // 关闭调试信息
		header := req.Header{
			"Accept": "application/json",
		}

		baseTransferUrl := config.Get("internal.internal_game_service") + config.Get("internal_api.balance_url")
		TransferUrl := baseTransferUrl + "?platform=" + platform + "&user_id=" + strconv.Itoa(userId)
		result := []UserGameAccount{}
		dbSession := common.Mysql(platform)
		defer dbSession.Close()

		for _, v := range gameVenues {
			if v.IsOnline != 1 { // 如果钱包已经下线则跳过
				continue
			}

			param := req.Param{
				"game_code": v.Code,
			}
			gc := UserGameAccount{ // 默认用户信息
				Code:    v.Code,
				Name:    v.Name,
				Account: "0.00",
			}
			r, err := req.Post(TransferUrl, header, param)
			if err != nil {
				log.Logger.Error(err.Error())
				result = append(result, gc)
				continue
			}
			transferResp := GameBalance{}
			err = r.ToJSON(&transferResp)
			if err != nil {
				fmt.Println("error: ", err)
				result = append(result, gc)
				continue
			}

			if v.Code == "CENTERWALLET" {
				sql := "select balance,username from accounts where user_id = " + strconv.Itoa(userId)
				res, err := dbSession.QueryString(sql)
				if err != nil {
					log.Err(err.Error())
					return result
				}
				if len(res) > 0 {
					gc.Account = res[0]["balance"]
				}
			}
			if transferResp.Data.GameCode == v.Code {
				gc.Account = strconv.Itoa(int(transferResp.Data.Balance)) // 不保留小数
			}
			if gc.Account == "" || gc.Account == "0" {
				gc.Account = "0.00"
			}
			result = append(result, gc)
		}

		return result
	},
}
