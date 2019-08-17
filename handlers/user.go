package handlers

import (
	"gin-api/helpers"
	"gin-api/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UserIndex(c *gin.Context) {
	users := models.UserGetList()
	page, _ := strconv.Atoi(c.Query("page"))
	start, end, currentPage, _ := helpers.Pager(page, len(users))

	c.JSON(http.StatusOK, gin.H{
		"data":    users[start:end],
		"total":   len(users),
		"current": currentPage,
	})
}

func UserShow(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	user := models.UserGetById(id)
	if user.ID == 0 {
		ReturnErrMsg(c, MSG_NOT_FOUND, nil)
		return
	}

	ReturnData(c, user)
}

func UserCreate(c *gin.Context) {
	var user models.User

	err := c.Bind(&user)
	if err != nil {
		ReturnErrMsg(c, MSG_FAILED, err)
		return
	}

	msg := user.Validate()
	if msg != "" {
		ReturnErrMsg(c, msg, nil)
		return
	}

	insertId := user.Create()
	if insertId == 0 {
		ReturnErrMsg(c, MSG_FAILED, err)
		return
	}

	ReturnMsg(c, MSG_SUCCESS)
}

func UserUpdate(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	user := models.UserGetById(id)
	if user.ID == 0 {
		ReturnErrMsg(c, MSG_NOT_FOUND, nil)
		return
	}

	err := c.Bind(&user)
	if err != nil {
		ReturnErrMsg(c, MSG_FAILED, err)
		return
	}

	msg := user.Validate()
	if msg != "" {
		ReturnErrMsg(c, msg, err)
		return
	}

	_ = user.Update()

	ReturnMsg(c, MSG_SUCCESS)
}

func UserDelete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	res := models.UserDeleteById(id)
	if res == false {
		ReturnErrMsg(c, MSG_FAILED, nil)
		return
	}

	ReturnMsg(c, MSG_SUCCESS)
}
