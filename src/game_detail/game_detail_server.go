package game_detail

import (
	"sports-admin/game_detail/ag"
	"sports-admin/game_detail/alb"
	"sports-admin/game_detail/avia"
	"sports-admin/game_detail/baoli"
	"sports-admin/game_detail/bbin"
	"sports-admin/game_detail/bg"
	"sports-admin/game_detail/bti"
	"sports-admin/game_detail/cq9"
	"sports-admin/game_detail/dt"
	"sports-admin/game_detail/ebet"
	"sports-admin/game_detail/hl"
	"sports-admin/game_detail/im"
	"sports-admin/game_detail/ime"
	"sports-admin/game_detail/jdb"
	"sports-admin/game_detail/ky"
	"sports-admin/game_detail/leg"
	"sports-admin/game_detail/leihuo"
	"sports-admin/game_detail/mg"
	"sports-admin/game_detail/ob"
	"sports-admin/game_detail/og"
	"sports-admin/game_detail/pt"
	"sports-admin/game_detail/saba"
	"sports-admin/game_detail/sg"
	"sports-admin/game_detail/ug"
	"sports-admin/game_detail/vg"
	"sports-admin/game_detail/vr"
	"sports-admin/game_detail/we"
	"sports-admin/game_detail/wm"
	"sports-admin/game_detail/xj188"
	models "sports-models"
)

type GameDetail interface {
	Data(billNo string, platform string) (*models.GameDetail, []*models.GameDetail)
}

type GameDetailSever struct {
	Detail GameDetail
}

var GameDetailList = map[string]GameDetail{
	"AG":     ag.NewGameDetail(),
	"ALB":    alb.NewGameDetail(),
	"BBIN":   bbin.NewGameDetail(),
	"BTI":    bti.NewGameDetail(),
	"DT":     dt.NewGameDetail(),
	"EBET":   ebet.NewGameDetail(),
	"HL":     hl.NewGameDetail(),
	"IM":     im.NewGameDetail(),
	"IME":    ime.NewGameDetail(),
	"JDB":    jdb.NewGameDetail(),
	"KY":     ky.NewGameDetail(),
	"LEG":    leg.NewGameDetail(),
	"LEIHUO": leihuo.NewGameDetail(),
	"MG":     mg.NewGameDetail(),
	"PT":     pt.NewGameDetail(),
	"SABA":   saba.NewGameDetail(),
	"VR":     vr.NewGameDetail(),
	"VG":     vg.NewGameDetail(),
	"BG":     bg.NewGameDetail(),
	"CQ9":    cq9.NewGameDetail(),
	"AVIA":   avia.NewGameDetail(),
	"SG":     sg.NewGameDetail(),
	"XJ188":  xj188.NewGameDetail(),
	"BAOLI":  baoli.NewGameDetail(),
	"WE":     we.NewGameDetail(),
	"UG":     ug.NewGameDetail(),
	"OG":     og.NewGameDetail(),
	"WM":     wm.NewGameDetail(),
	"OB":     ob.NewGameDetail(),
}

func NewGameDetail(gameCode string) GameDetail {
	temp := GameDetailList[gameCode]
	return temp
}
