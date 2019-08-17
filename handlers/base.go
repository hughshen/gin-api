package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

const (
	MSG_FAILED    string = "Failed"
	MSG_SUCCESS   string = "Successfully"
	MSG_NOT_FOUND string = "Not found"
)

func ShowErrLog(err error) {
	log.Println(err)
}

func ReturnMsg(c *gin.Context, message string) {
	c.JSON(http.StatusOK, gin.H{
		"error":   0,
		"message": message,
	})
}

func ReturnErrMsg(c *gin.Context, message string, err error) {
	if err != nil {
		ShowErrLog(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   1,
		"message": message,
	})
}

func ReturnData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}

func ReturnNotFound(c *gin.Context, err error) {
	if err != nil {
		ShowErrLog(err)
	}

	c.HTML(http.StatusNotFound, "404.html", gin.H{})
}
