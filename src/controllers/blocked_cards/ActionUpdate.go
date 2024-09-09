package blocked_cards

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionUpdate = &base_controller.ActionUpdate{
	Model: models.BlockedCards,
	Row: func() interface{} {
		return &models.BlockedCard{}
	},
	ViewFile: "blacklists/cards_edit.html",
}
