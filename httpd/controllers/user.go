package controllers

import (
	"go-demo/httpd/models"
	"go-demo/httpd/services"
	"go-demo/httpd/utils"
	"go-demo/internal/logger"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddUser(c *gin.Context) {
	var m models.User
	if err := c.ShouldBindJSON(&m); err != nil {
		utils.ResponseFailureValidatorParameter(c, err)
		return
	}
	_, err := services.AddUser(&m)
	if err != nil {
		logger.GetLogger().Error().Msgf("Add alarm user error %s", err.Error())
		utils.ResponseFailureServer(c, err)
		return
	}
	utils.ResponseSuccess(c, gin.H{"id": m.Id})
}

func GetUserPage(c *gin.Context) {
	obj, isExist := c.GetQuery("pageIndex")
	if isExist != true {
		utils.ResponseFailureParameter(c, "参数pageIndex不能为空")
		return
	}
	pageIndex, err := strconv.Atoi(obj)
	if err != nil {
		utils.ResponseFailureParameter(c, "参数pageIndex必须是整数")
		return
	}
	obj, isExist = c.GetQuery("pageSize")
	if isExist != true {
		utils.ResponseFailureParameter(c, "参数pageSize不能为空")
		return
	}
	pageSize, err := strconv.Atoi(obj)
	if err != nil {
		utils.ResponseFailureParameter(c, "参数pageSize必须是整数")
		return
	}
	name := c.Query("name")
	data, err := services.GetUserPage(pageIndex, pageSize, name)
	if err != nil {
		utils.ResponseFailureParameter(c, err.Error())
		return
	}
	utils.ResponseSuccess(c, data)
}
