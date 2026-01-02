package utils

import (
	"encoding/json"
	"fmt"
	myvalidator "go-demo/internal/validator"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func responseGetErrorMsg(err interface{}) string {
	errorMsg := ""
	switch v := err.(type) {
	case error:
		errorMsg = v.Error()
	case string:
		errorMsg = v
	default:
		errorMsg = fmt.Sprintf("错误信息： %v", v)
	}
	return errorMsg
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	if data == nil {
		c.JSON(200, gin.H{})
	} else {
		c.JSON(200, data)
	}
}

func ResponseBadRequest(c *gin.Context, httpCode int, errorCode int, err interface{}, errorData interface{}) {
	errorMsg := responseGetErrorMsg(err)
	data := gin.H{"errorCode": errorCode, "errorMsg": errorMsg, "errorData": errorData}
	c.JSON(httpCode, data)
}

func ResponseErrorValidatorParameter(c *gin.Context, err error) {
	var errorMsg string
	errorData := make(map[string]string)
	switch errs := err.(type) {
	case nil:
		errorMsg = "error is nil"
	case validator.ValidationErrors:
		translator := myvalidator.GetTranslator()
		errorData = myvalidator.RemoveStructName(errs.Translate(translator))
		for _, v := range errorData {
			if errorMsg == "" {
				errorMsg = v
			} else {
				errorMsg = fmt.Sprintf("%s: %s", errorMsg, v)
			}
		}
	case *json.UnmarshalTypeError:
		errorData[errs.Field] = fmt.Sprintf("类型错误: 期望类型 %s", errs.Type.String())
		errorMsg = fmt.Sprintf("%s 类型错误，期望类型 %s", errs.Field, errs.Type.String())
	default:
		errorMsg = err.Error()
	}
	ResponseBadRequest(c, 400, 54000, errorMsg, errorData)
}

func ResponseErrorParameter(c *gin.Context, err interface{}) {
	ResponseBadRequest(c, 400, 54000, err, nil)
}

func ResponseMsgUnauthorized(c *gin.Context, err interface{}, errorData string) {
	ResponseBadRequest(c, 401, 54010, err, errorData)
}

func ResponseErrorForbidden(c *gin.Context, err interface{}, errorData string) {
	ResponseBadRequest(c, 403, 54030, err, errorData)
}

func ResponseErrorNotFoundData(c *gin.Context, err interface{}) {
	ResponseBadRequest(c, 404, 54040, err, nil)
}

func ResponseErrorNotFoundPage(c *gin.Context, err interface{}) {
	ResponseBadRequest(c, 404, 54041, err, nil)
}

func ResponseErrorServer(c *gin.Context, err interface{}) {
	ResponseBadRequest(c, 500, 5400, err, nil)
}
