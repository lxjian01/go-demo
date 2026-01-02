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
	c.JSON(200, gin.H{"code": 0, "message": "success", "data": data})
}

func ResponseFailure(c *gin.Context, httpCode int, errorCode int, err interface{}, errorData interface{}) {
	errorMsg := responseGetErrorMsg(err)
	c.JSON(httpCode, gin.H{"code": errorCode, "message": errorMsg, "data": errorData})
}

func ResponseFailureValidatorParameter(c *gin.Context, err error) {
	var errorMsg string
	errorData := make(map[string]string)
	switch errs := err.(type) {
	case nil:
		errorMsg = "error is nil"
	case validator.ValidationErrors:
		errorData = myvalidator.Translate(errs)
		for _, v := range errorData {
			if errorMsg == "" {
				errorMsg = v
			} else {
				errorMsg = errorMsg + "; " + v
			}
		}
	case *json.UnmarshalTypeError:
		// errs.Field 是 Go 结构体字段名
		jsonField := myvalidator.ConvertFieldToJSONTag(errs.Struct, errs.Field)
		errorData[jsonField] = fmt.Sprintf("类型错误: 期望类型 %s", errs.Type.String())
		errorMsg = fmt.Sprintf("%s 类型错误，期望类型 %s", jsonField, errs.Type.String())
	default:
		errorMsg = err.Error()
	}
	ResponseFailure(c, 400, 54000, errorMsg, errorData)
}

func ResponseFailureParameter(c *gin.Context, err interface{}) {
	ResponseFailure(c, 400, 54000, err, nil)
}

func ResponseFailureUnauthorized(c *gin.Context, err interface{}) {
	ResponseFailure(c, 401, 54010, err, nil)
}

func ResponseFailureForbidden(c *gin.Context, err interface{}) {
	ResponseFailure(c, 403, 54030, err, nil)
}

func ResponseFailureNotFoundData(c *gin.Context, err interface{}) {
	ResponseFailure(c, 404, 54040, err, nil)
}

func ResponseFailureNotFoundPage(c *gin.Context, err interface{}) {
	ResponseFailure(c, 404, 54041, err, nil)
}

func ResponseFailureServer(c *gin.Context, err interface{}) {
	ResponseFailure(c, 500, 5400, err, nil)
}
