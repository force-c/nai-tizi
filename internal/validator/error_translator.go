package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// fieldNameMap 字段名称映射（英文 -> 中文）
var fieldNameMap = map[string]string{
	"OrgId":       "组织ID",
	"OrgIds":      "组织ID列表",
	"ParentId":    "父组织ID",
	"OrgName":     "组织名称",
	"OrgCode":     "组织编码",
	"OrgType":     "组织类型",
	"Leader":      "负责人",
	"Phone":       "联系电话",
	"SortOrder":   "显示顺序",
	"UserId":      "用户ID",
	"UserName":    "用户名",
	"NickName":    "昵称",
	"Password":    "密码",
	"UserType":    "用户类型",
	"Email":       "邮箱",
	"Phonenumber": "手机号",
	"Sex":         "性别",
	"Avatar":      "头像",
	"Status":      "状态",
	"Remark":      "备注",
	"NewPassword": "新密码",
	"UserIds":     "用户ID列表",
}

// TranslateValidationError 翻译验证错误为友好的中文提示
func TranslateValidationError(err error) string {
	if err == nil {
		return ""
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}

	var messages []string
	for _, e := range validationErrors {
		messages = append(messages, translateSingleError(e))
	}

	return strings.Join(messages, "; ")
}

// translateSingleError 翻译单个验证错误
func translateSingleError(e validator.FieldError) string {
	fieldName := getFieldName(e.Field())

	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s不能为空", fieldName)
	case "min":
		if e.Type().Kind() == 24 { // string
			return fmt.Sprintf("%s长度不能少于%s个字符", fieldName, e.Param())
		}
		return fmt.Sprintf("%s不能小于%s", fieldName, e.Param())
	case "max":
		if e.Type().Kind() == 24 { // string
			return fmt.Sprintf("%s长度不能超过%s个字符", fieldName, e.Param())
		}
		return fmt.Sprintf("%s不能大于%s", fieldName, e.Param())
	case "len":
		return fmt.Sprintf("%s长度必须为%s位", fieldName, e.Param())
	case "email":
		return fmt.Sprintf("%s格式不正确", fieldName)
	case "oneof":
		return fmt.Sprintf("%s的值必须是以下之一: %s", fieldName, e.Param())
	case "cnphone":
		return fmt.Sprintf("%s格式不正确，请输入11位手机号", fieldName)
	case "mac":
		return fmt.Sprintf("%s格式不正确", fieldName)
	case "sn":
		return fmt.Sprintf("%s格式不正确", fieldName)
	case "dive":
		return fmt.Sprintf("%s中的数据验证失败", fieldName)
	default:
		return fmt.Sprintf("%s验证失败: %s", fieldName, e.Tag())
	}
}

// getFieldName 获取字段的中文名称
func getFieldName(field string) string {
	if name, ok := fieldNameMap[field]; ok {
		return name
	}
	return field
}
