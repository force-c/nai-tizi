package response

import "github.com/gin-gonic/gin"

const (
	CodeOK              = 200
	CodeBadRequest      = 400
	CodeUnauthorized    = 401
	CodeForbidden       = 403
	CodeNotFound        = 404
	CodeTimeout         = 408
	CodeTooManyRequests = 429
	CodeServerError     = 500
	CodeInvalidParam    = 400
)

// Response 统一响应结构
//
//	@Description	API 统一响应格式
type Response struct {
	Code int         `json:"code" example:"200"`        // 业务状态码
	Msg  string      `json:"msg" example:"success"`     // 响应消息
	Data interface{} `json:"data" swaggertype:"object"` // 响应数据
}

// PageResponse 分页响应结构
//
//	@Description	分页查询统一响应格式
type PageResponse struct {
	Code  int         `json:"code" example:"200"`    // 业务状态码
	Msg   string      `json:"msg" example:"success"` // 响应消息
	Rows  interface{} `json:"rows"`                  // 数据列表
	Total int64       `json:"total" example:"100"`   // 总记录数
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{Code: CodeOK, Msg: "success", Data: data})
}

func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(200, Response{Code: CodeOK, Msg: msg, Data: data})
}

func Fail(c *gin.Context, msg string)        { c.JSON(200, Response{Code: CodeServerError, Msg: msg}) }
func FailWithMsg(c *gin.Context, msg string) { Fail(c, msg) }

func BadRequest(c *gin.Context, msg string) {
	c.JSON(200, Response{Code: CodeBadRequest, Msg: msg})
}

func Unauthorized(c *gin.Context, msg string) {
	c.JSON(200, Response{Code: CodeUnauthorized, Msg: msg})
}

func Forbidden(c *gin.Context, msg string) {
	c.JSON(200, Response{Code: CodeForbidden, Msg: msg})
}

func NotFound(c *gin.Context, msg string) {
	c.JSON(200, Response{Code: CodeNotFound, Msg: msg})
}

func InternalServerError(c *gin.Context, msg string) {
	c.JSON(200, Response{Code: CodeServerError, Msg: msg})
}

func PageSuccess(c *gin.Context, rows interface{}, total int64) {
	c.JSON(200, PageResponse{Code: CodeOK, Msg: "success", Rows: rows, Total: total})
}

func SuccessCode(c *gin.Context, code int, data interface{}) {
	c.JSON(200, Response{Code: code, Msg: "success", Data: data})
}

func FailCode(c *gin.Context, code int, msg string) { c.JSON(200, Response{Code: code, Msg: msg}) }
