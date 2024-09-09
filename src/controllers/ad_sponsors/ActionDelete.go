package ad_sponsors

import (
	"sports-admin/controllers/base_controller"
	models "sports-models"
)

var ActionDelete = &base_controller.ActionDelete{
	Model: models.AdSponsors,
}
