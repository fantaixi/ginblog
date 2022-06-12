package validator

import (
	"fmt"
	"ginblog/utils/errmsg"
	"github.com/go-playground/locales/zh_Hans_CN"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTrans "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

func Validate(data interface{}) (string, int) {
	validate := validator.New()
	uni := ut.New(zh_Hans_CN.New())
	trans,_ := uni.GetTranslator("zh_Hans_CN")
	err := zhTrans.RegisterDefaultTranslations(validate,trans)
	if err != nil {
		fmt.Println("err:",err)
	}
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		lable := field.Tag.Get("label")
		return lable
	})
	//这里能保证传进来的都是结构体，所以不用断言，如果是其他数据，就需要断言
	err = validate.Struct(data)
	if err != nil {
		for _, v := range err.(validator.ValidationErrors) {
			return v.Translate(trans),errmsg.ERROR
		}
	}
	return "",errmsg.SUCCESS
}

