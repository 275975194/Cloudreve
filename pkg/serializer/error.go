package serializer

import "github.com/gin-gonic/gin"

// Response 基础序列化器
type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data,omitempty"`
	Msg   string      `json:"msg"`
	Error string      `json:"error,omitempty"`
}

// AppError 应用错误，实现了error接口
type AppError struct {
	Code     int
	Msg      string
	RawError error
}

// NewError 返回新的错误对象 todo:测试 还有下面的
func NewError(code int, msg string, err error) AppError {
	return AppError{
		Code:     code,
		Msg:      msg,
		RawError: err,
	}
}

// WithError 将应用error携带标准库中的error
func (err *AppError) WithError(raw error) AppError {
	err.RawError = raw
	return *err
}

// Error 返回业务代码确定的可读错误信息
func (err AppError) Error() string {
	return err.Msg
}

// 三位数错误编码为复用http原本含义
// 五位数错误编码为应用自定义错误
// 五开头的五位数错误编码为服务器端错误，比如数据库操作失败
// 四开头的五位数错误编码为客户端错误，有时候是客户端代码写错了，有时候是用户操作错误
const (
	// CodeCheckLogin 未登录
	CodeCheckLogin = 401
	// CodeNoRightErr 未授权访问
	CodeNoRightErr = 403
	// CodeUploadFailed 上传出错
	CodeUploadFailed = 40002
	// CodeCreateFolderFailed 目录创建失败
	CodeCreateFolderFailed = 40003
	// CodeObjectExist 对象已存在
	CodeObjectExist = 40004
	// CodeDBError 数据库操作失败
	CodeDBError = 50001
	// CodeEncryptError 加密失败
	CodeEncryptError = 50002
	// CodePolicyNotAllowed 当前存储策略不允许
	CodePolicyNotAllowed = 50003
	// CodeIOFailed IO操作失败
	CodeIOFailed = 50004
	//CodeParamErr 各种奇奇怪怪的参数错误
	CodeParamErr = 40001
	// CodeNotSet 未定错误，后续尝试从error中获取
	CodeNotSet = -1
)

// DBErr 数据库操作失败
func DBErr(msg string, err error) Response {
	if msg == "" {
		msg = "数据库操作失败"
	}
	return Err(CodeDBError, msg, err)
}

// ParamErr 各种参数错误
func ParamErr(msg string, err error) Response {
	if msg == "" {
		msg = "参数错误"
	}
	return Err(CodeParamErr, msg, err)
}

// Err 通用错误处理
func Err(errCode int, msg string, err error) Response {
	// 如果错误code未定，则尝试从AppError中获取
	if errCode == CodeNotSet {
		if appError, ok := err.(AppError); ok {
			errCode = appError.Code
			err = appError.RawError
		}
	}
	res := Response{
		Code: errCode,
		Msg:  msg,
	}
	// 生产环境隐藏底层报错
	if err != nil && gin.Mode() != gin.ReleaseMode {
		res.Error = err.Error()
	}
	return res
}
