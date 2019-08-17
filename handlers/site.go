package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SiteNotFound(c *gin.Context) {
	ReturnNotFound(c, nil)
}

func SiteHome(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		name = "world"
	}

	c.HTML(http.StatusOK, "home.html", gin.H{
		"name": name,
	})
}
