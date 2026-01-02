package utils

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct{}

func (resp *Response) ToSuccess(c *gin.Context, data interface{}) {
	if data == nil {
		c.JSON(200, gin.H{})
	} else {
		c.JSON(200, data)
	}
}

func (resp *Response) ToBadRequest(c *gin.Context, httpCode int, errorCode int, err interface{}, errorData interface{}) {
	errorMsg := resp.getErrorMsg(err)
	c.Set(errorMsg, errorMsg)
	data := gin.H{"errorCode": errorCode, "errorMsg": errorMsg, "errorData": errorData}
	c.JSON(httpCode, data)
}

func (resp *Response) getErrorMsg(err interface{}) string {
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

func (resp *Response) ToErrorValidatorParameter(c *gin.Context, err error) {
	var errorMsg string
	errorData := make(map[string]string)
	switch errs := err.(type) {
	case nil:
		errorMsg = "error is nil"
	case validator.ValidationErrors:
		translator := GetTranslator()
		errorData = RemoveStructName(errs.Translate(translator))
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
	resp.ToBadRequest(c, 400, 54000, errorMsg, errorData)
}

func (resp *Response) ToErrorParameter(c *gin.Context, err interface{}) {
	resp.ToBadRequest(c, 400, 54000, err, nil)
}

func (resp *Response) ToMsgUnauthorized(c *gin.Context, err interface{}, errorData string) {
	resp.ToBadRequest(c, 401, 54010, err, errorData)
}

func (resp *Response) ToErrorForbidden(c *gin.Context, err interface{}, errorData string) {
	resp.ToBadRequest(c, 403, 54030, err, errorData)
}

func (resp *Response) ToErrorNotFoundData(c *gin.Context, err interface{}) {
	resp.ToBadRequest(c, 404, 54040, err, nil)
}

func (resp *Response) ToErrorNotFoundPage(c *gin.Context, err interface{}) {
	resp.ToBadRequest(c, 404, 54041, err, nil)
}

func (resp *Response) ToErrorServer(c *gin.Context, err interface{}) {
	resp.ToBadRequest(c, 500, 5400, err, nil)
}
