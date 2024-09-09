package controllers

import (
	"sports-admin/controllers/ad_carousels"
)

// AdCarousels 轮播图
var AdCarousels = ad_carousels.AdCarousels{
	ActionList:   ad_carousels.ActionList,
	ActionCreate: ad_carousels.ActionCreate,
	ActionUpdate: ad_carousels.ActionUpdate,
	ActionSave:   ad_carousels.ActionSave,
	ActionDelete: ad_carousels.ActionDelete,
	ActionState:  ad_carousels.ActionState,
}
