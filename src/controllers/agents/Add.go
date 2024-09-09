package agents

import (
	"sports-common/response"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func (ths *Agents) Add(c *gin.Context) {
	name := c.Query("username")
	id := c.Query("id")
	response.Render(c, "agents/agents_useradd.html", pongo2.Context{"username": name, "id": id})
}
