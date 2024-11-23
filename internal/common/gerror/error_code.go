package gerror

// ErrorMap 全局错误码 映射
var ErrorMap = make(map[string]*ErrorCode)

type ErrorCode struct {
	Code    string
	Message string
	Parent  *ErrorCode
	Alarm   bool
}

func NewErrorCode(code string, message string, parent *ErrorCode, alarm bool) *ErrorCode {
	if len(code) == 0 {
		return nil
	}
	e := &ErrorCode{
		Code:    code,
		Message: message,
		Parent:  parent,
		Alarm:   alarm,
	}
	ErrorMap[code] = e
	return e
}

func GetErrorCode(code string) *ErrorCode {
	return ErrorMap[code]
}

// 系统通用错误-错误码定义
var (
	SystemError        = NewErrorCode("SystemError", "Something went wrong with the system. Please try again later", nil, true) // 系统错误
	RemoteGrpcError    = NewErrorCode("RemoteGrpcError", "Remote service failure", nil, true)                                   // 远程grpc错误
	RemoteHttpError    = NewErrorCode("RemoteHttpError", "Remote http service failure", nil, true)                              // 远程http错误
	ParamError         = NewErrorCode("ParamError", "Param illegal:%s", nil, false)
	DataNotExist       = NewErrorCode("DataNotExist", "Data not exist", nil, false)                 // 数据不存在
	NoPrivilege        = NewErrorCode("NoPrivilege", "No privilege", nil, false)                    // 数据访问没有权限
	UserNotLogin       = NewErrorCode("UserNotLogin", "User not login", nil, false)                 // 用户未登录
	ServiceTimeOut     = NewErrorCode("ServiceTimeOut", "Service time out", nil, true)              // 超时
	RedisError         = NewErrorCode("RedisError", "cache gerror", nil, true)                      // redis错误
	DataConvertError   = NewErrorCode("DataConvertError", "Parameter conversion gerror", nil, true) // 数据转换失败
	DbError            = NewErrorCode("DbError", "DB gerror", nil, true)                            // DB错误
	UrlCanNotNull      = NewErrorCode("UrlCanNotNull", "Url can not null", nil, true)
	UrlProtocolMissing = NewErrorCode("UrlProtocolMissing", "Url protocol missing", nil, true)
	TooManyRequest     = NewErrorCode("TooManyRequest", "Too many request", nil, true)
)

func ErrSystemError() *GnBizError {
	return NewBizError(SystemError)
}

func ErrRemoteGrpcError() *GnBizError {
	return NewBizError(RemoteGrpcError)
}

func ErrRemoteHttpError() *GnBizError {
	return NewBizError(RemoteHttpError)
}

func ErrParamError() *GnBizError {
	return NewBizError(ParamError)
}

func ErrDataNotExist() *GnBizError {
	return NewBizError(DataNotExist)
}

func ErrNoPrivilege() *GnBizError {
	return NewBizError(NoPrivilege)
}

func ErrUserNotLogin() *GnBizError {
	return NewBizError(UserNotLogin)
}

func ErrServiceTimeOut() *GnBizError {
	return NewBizError(ServiceTimeOut)
}

func ErrRedisError() *GnBizError {
	return NewBizError(RedisError)
}

func ErrDataConvertError() *GnBizError {
	return NewBizError(DataConvertError)
}

func ErrDbError() *GnBizError {
	return NewBizError(DbError)
}

func ErrUrlCanNotNull() *GnBizError {
	return NewBizError(UrlCanNotNull)
}

func ErrUrlProtocolMissing() *GnBizError {
	return NewBizError(UrlProtocolMissing)
}
