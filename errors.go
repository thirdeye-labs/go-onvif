package onvif

import "fmt"

type OnvifErr struct {
	subCode string
	Detail  string
}

type ErrOperationProhibited OnvifErr

func NewErrOperationProhibited(detail string) ErrOperationProhibited {
	return ErrOperationProhibited{
		subCode: "OperationProhibited",
		Detail:  detail,
	}
}
func (e ErrOperationProhibited) Error() string {
	return fmt.Sprintf("%s: %s", e.subCode, e.Detail)
}

type ErrNewUnsupportedError OnvifErr

func NewUnsupportedError(subCode, detail string) ErrNewUnsupportedError {
	return ErrNewUnsupportedError{
		subCode: subCode,
		Detail:  detail,
	}
}
func (e ErrNewUnsupportedError) Error() string {
	return fmt.Sprintf("%s: %s", e.subCode, e.Detail)
}
