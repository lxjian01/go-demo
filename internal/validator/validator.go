package validator

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var (
	trans ut.Translator
	once  sync.Once
)

func GetTranslator() ut.Translator {
	return trans
}

func InitTrans(locale string) error {
	var initErr error
	once.Do(func() {
		initErr = initTranslator(locale)
	})
	return initErr
}

func initTranslator(locale string) (err error) {
	//修改gin框架中Validator引擎属性，实现自定制
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
		zhT := zh.New() //中文翻译器
		enT := en.New() //英文翻译器
		// 第一个参数是备用(fallback)语言环境
		// 后面参数是应该支持语言环境(可支持多个)
		uni := ut.New(enT, zhT, enT)
		// locale通常取决于http请求'Accept-language'
		var ok bool
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("uni.GetTranslator(%s) failed", locale)
		}
		// 注册翻译器
		switch locale {
		case "en":
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		case "zh":
			err = zhTranslations.RegisterDefaultTranslations(v, trans)
		default:
			err = enTranslations.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

// ConvertFieldToJSONTag 根据结构体类型和 Go 字段名，返回 JSON tag
func ConvertFieldToJSONTag(structType interface{}, goField string) string {
	t := reflect.TypeOf(structType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if field, ok := t.FieldByName(goField); ok {
		tag := field.Tag.Get("json")
		if tag == "" || tag == "-" {
			return goField
		}
		return strings.Split(tag, ",")[0] // 去掉 `omitempty` 等选项
	}
	return goField // fallback
}

// Translate 接收 validator.ValidationErrors，返回 key=字段名, value=翻译后的错误信息
func Translate(errs validator.ValidationErrors) map[string]string {
	if errs == nil {
		return nil
	}

	errorData := make(map[string]string, len(errs))
	for _, e := range errs {
		field := e.StructField() // 结构体字段名
		if trans != nil {
			errorData[field] = e.Translate(trans)
		} else {
			errorData[field] = e.Error()
		}
	}
	return RemoveStructName(errorData) // 去掉结构体前缀
}

func RemoveStructName(fields map[string]string) map[string]string {
	result := make(map[string]string, len(fields))

	for field, err := range fields {
		if idx := strings.Index(field, "."); idx != -1 {
			result[field[idx+1:]] = err
		} else {
			result[field] = err
		}
	}
	return result
}
