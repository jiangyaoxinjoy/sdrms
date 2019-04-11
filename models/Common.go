package models

import "sdrms/enums"

// JsonResult 用于返回ajax请求的基类
type JsonResult struct {
	Code enums.JsonResultCode `json:"code"`
	Msg  string               `json:"msg"`
	Obj  interface{}          `json:"obj"`
}

type JsonResultEditormd struct {
	Success enums.JsonResultCode `json:"success"`
	Message string               `json:"message"`
	Url     string               `json:"url"`
}

// BaseQueryParam 用于查询的类
type BaseQueryParam struct {
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Offset int64  `json:"offset"`
	Limit  int    `json:"limit"`
}
