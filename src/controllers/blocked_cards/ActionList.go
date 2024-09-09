package blocked_cards

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionList = &base_controller.ActionList{
	Model:    models.BlockedCards,
	ViewFile: "blacklists/cards.html",
	Rows: func() interface{} {
		return &[]models.BlockedCard{}
	},
	QueryCond: map[string]interface{}{
		"card_number": "%",
		"user_names":  "%",
	},
}
