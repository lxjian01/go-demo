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
	var resp utils.Response
	var m models.User
	if err := c.ShouldBindJSON(&m); err != nil {
		resp.ToError(c, err)
		return
	}
	_, err := services.AddUser(&m)
	if err != nil {
		logger.GetLogger().Error().Msgf("Add alarm user error %s", err.Error())
		resp.ToError(c, err)
		return
	}
	resp.Data = gin.H{"id": m.Id}
	resp.ToSuccess(c)
}

func GetUserPage(c *gin.Context) {
	resp := &utils.Response{}
	obj, isExist := c.GetQuery("pageIndex")
	if isExist != true {
		resp.ToMsgBadRequest(c, "参数pageIndex不能为空")
		return
	}
	pageIndex, err := strconv.Atoi(obj)
	if err != nil {
		resp.ToMsgBadRequest(c, "参数pageIndex必须是整数")
		return
	}
	obj, isExist = c.GetQuery("pageSize")
	if isExist != true {
		resp.ToMsgBadRequest(c, "参数pageSize不能为空")
		return
	}
	pageSize, err := strconv.Atoi(obj)
	if err != nil {
		resp.ToMsgBadRequest(c, "参数pageSize必须是整数")
		return
	}
	name := c.Query("name")
	data, err := services.GetUserPage(pageIndex, pageSize, name)
	if err != nil {
		resp.ToMsgBadRequest(c, err.Error())
		return
	}

	resp.Data = data
	resp.ToSuccess(c)
}
