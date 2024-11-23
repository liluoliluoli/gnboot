package gerror

import (
	"fmt"
	"strings"
)

type (
	GnBizError struct {
		code       string
		message    string
		cause      error
		values     []interface{}
		stacktrace string
	}
)

func (e GnBizError) Values() []interface{} {
	return e.values
}

func (e *GnBizError) Error() string {
	sb := &strings.Builder{}
	sb.WriteString("<code: ")
	sb.WriteString(e.code)
	sb.WriteString(", message: ")
	sb.WriteString(e.message)

	if e.cause != nil {
		sb.WriteString(", cause: ")
		sb.WriteString(e.cause.Error())
	}

	sb.WriteString(">")
	return sb.String()
}

func (e *GnBizError) Code() string {
	return e.code
}

func (e *GnBizError) Message() string {
	//if e.message == "" {
	//	if e.cause != nil {
	//		return e.cause.Error()
	//	}
	//}

	return e.message
}

func (e *GnBizError) Cause() error {
	return e.cause
}

func (e *GnBizError) Stacktrace() string {
	return e.stacktrace
}

func (e *GnBizError) Detail() []byte {
	return []byte(e.stacktrace)
}

func newGnBizError(code, message string, cause error, values ...interface{}) *GnBizError {
	err := &GnBizError{
		code:       code,
		message:    message,
		cause:      cause,
		values:     values,
		stacktrace: dumpStacktrace(2).String(),
	}

	return err
}

func NewBizError(errorCode *ErrorCode, values ...interface{}) *GnBizError {
	return newGnBizError(errorCode.Code, errorCode.Message, nil, values...)
}

func NewBizErrorWithCause(errorCode *ErrorCode, cause error, values ...interface{}) *GnBizError {
	return newGnBizError(errorCode.Code, errorCode.Message, cause, values...)
}

// GetCodeAndMessage 包装获取err的code和msg, msg优先用err中获取，再从errorcode获取；优先从parent err获取
func GetCodeAndMessage(err error, isWrapper bool, i18nCode string) (string, string, bool) {
	switch err := err.(type) {
	case nil:
		return "", "", false
	case *GnBizError:
		code := err.Code()
		alarm := true
		ec, ok := ErrorMap[code]
		if ok {
			if ec.Parent != nil {
				code = ec.Parent.Code
				pec, ok := ErrorMap[code]
				if ok {
					ec = pec
				}
			}
		}
		if ec != nil {
			alarm = ec.Alarm
		}

		message := err.Message()
		if len(err.Message()) == 0 && ec != nil {
			message = ec.Message
		}

		if err.Values() != nil && len(err.Values()) > 0 {
			if strings.Contains(message, "%") {
				message = fmt.Sprintf(message, err.Values()...)
			} else {
				for _, value := range err.Values() {
					message += fmt.Sprintf(", %v", value)
				}
			}
		}
		return code, message, alarm

	default:
		if isWrapper {
			return GetCodeAndMessage(NewBizErrorWithCause(SystemError, err), true, i18nCode)
		} else {
			return err.Error(), err.Error(), true
		}
	}
}

// GetPrintErrStacks 打印cause err的堆栈, i18nCode默认为:en
func GetPrintErrStacks(err error, i18nCode string) []string {
	var errStacks []string
	switch err := err.(type) {
	case nil:

	case *GnBizError:
		if err.Cause() != nil {
			errStacks = GetPrintErrStacks(err.Cause(), i18nCode)
		}
		code, message, _ := GetCodeAndMessage(err, false, i18nCode)
		errStacks = append(errStacks, fmt.Sprintf("code: %s, message: %s, stacktrace: %s", code, message, err.stacktrace))
	default:
		errStacks = append(errStacks, fmt.Sprintf("stacktrace: %+v", err))
	}
	return errStacks
}

func GetPrintErrStackStr(err error, i18nCode string) []string {
	return GetPrintErrStacks(err, i18nCode)
}
