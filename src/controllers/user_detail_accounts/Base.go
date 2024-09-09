package user_detail_accounts

import models "sports-models"

// TransferState 转账信息
type TransferState struct {
	Errcode int    `json:"errcode"`
	Message string `json:"message"`
}

// RebateAndVip 返水VIP
type RebateAndVip struct {
	models.UserRebateRecord `xorm:"extends"`
	Vip                     int32
}

// SumAgent 统计
type SumAgent struct {
	Value float64 `json:"value"`
}

// UserDepositSumStruct 用户存款汇总
type UserDepositSumStruct struct {
	Money    float64 `json:"money"`
	Discount float64 `json:"discount"`
}

// GameAccount 游戏钱包
type GameAccount struct {
	PlayName string  `json:"playname" xorm:"playname"`
	GameCode string  `json:"game_code"`
	Money    float64 `json:"money"`
	Created  int     `json:"created"`
	Updated  int     `json:"updated"`
}

// UserDetailWin 用户输赢
type UserDetailWin struct {
	GameCode   string  `json:"game_code"`
	BetMoney   float64 `json:"bet_money"`
	NetMoney   float64 `json:"net_money"`
	ValidMoney float64 `json:"valid_money"`
	// GameType   uint8   `json:"game_type"`
	// ExtendStr  string  `json:"extend_str"`
	// Status     int8    `json:"status"`
	// CreatedAt  uint32  `json:"created_at"`
	// UpdatedAt  uint32  `json:"updated_at"`
}

// UserDetailAccounts 会员详情 - 账户信息
type UserDetailAccounts struct{}
