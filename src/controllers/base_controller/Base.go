package base_controller

// IdRecord 只有id字段的数据记录
type IdRecord struct {
	Id uint32 `json:"id"`
}

type SumStruct struct {
	Money float64
}

type TransferResp struct {
	Errcode int    `json:"errcode"`
	Message string `json:"message"`
	Data    struct {
		Money       string `json:"money"`
		GameBalance []struct {
			GameCode string `json:"game_code"`
			Balance  string `json:"balance"`
		} `json:"game_balance"`
	} `json:"data"`
}

type BalanceResp struct {
	Errcode int    `json:"errcode"`
	Message string `json:"message"`
	Data    struct {
		Balance  float64 `json:"balance"`
		GameCode string  `json:"game_code"`
		State    int     `json:"state"`
	} `json:"data"`
}
