package blocked_cards

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionCreate = &base_controller.ActionCreate{
	Model:    models.BlockedCards,
	ViewFile: "blacklists/cards_edit.html",
}
