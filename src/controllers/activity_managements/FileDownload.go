package activity_managements

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (ths *ActivityManagements) FileDownload(c *gin.Context) {
	filename := "dividend.xlsx"
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Writer.Header().Add("Content-Type", "application/octet-stream")
	c.File("./uploads/Excel/" + filename)
}
