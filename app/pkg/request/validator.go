package request

import (
   "github.com/go-playground/validator/v10"
   cErr "github.com/jassue/gin-wire/app/pkg/error"
)

type Validator interface {
   GetMessages() ValidatorMessages
}

type ValidatorMessages map[string]string

// GetError 获取验证错误
func GetError(request interface{}, err error) *cErr.Error {
   if _, isValidatorErrors := err.(validator.ValidationErrors); isValidatorErrors {
      _, isValidator := request.(Validator)

      for _, v := range err.(validator.ValidationErrors) {
         // 若 request 结构体实现 Validator 接口即可实现自定义错误信息
         if isValidator {
            if message, exist := request.(Validator).GetMessages()[v.Field() + "." + v.Tag()]; exist {
               return cErr.ValidateErr(message)
            }
         }
         return cErr.ValidateErr(v.Error())
      }
   }

   return cErr.ValidateErr("Parameter error")
}
