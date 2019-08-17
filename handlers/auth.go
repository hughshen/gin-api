package handlers

import (
	"gin-api/models"
	"github.com/gin-gonic/gin"
)

func AuthLogin(c *gin.Context) {
	var form models.LoginForm

	c.Bind(&form)

	user, msg := form.Validate()
	if msg != "" {
		ReturnErrMsg(c, msg, nil)
		return
	}

	ReturnData(c, user)
}

func AuthHello(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	ReturnMsg(c, "hello, "+user.Username)
}

func AuthLogout(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	res := models.UserDeleteToken(user.ID)
	if res == true {
		ReturnMsg(c, MSG_SUCCESS)
	} else {
		ReturnErrMsg(c, MSG_FAILED, nil)
	}
}
