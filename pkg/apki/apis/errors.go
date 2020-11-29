package apis

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Err struct {
	Message string      `json:"message"`
	Code    string      `json:"code,omitempty"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Reason  string      `json:"reason,omitempty"`
	Cause   error       `json:"-"`
}

func (e Err) Error() string {
	return e.Message
}
func FromError(err interface{}) *Err {
	if err == nil {
		return nil
	}
	var e *Err
	switch v := err.(type) {
	case Err:
		e = &v
	case *Err:
		e = v
	case error:
		e = &Err{
			Message: v.Error(),
			Status:  http.StatusInternalServerError,
			Cause:   v,
		}
	case string:
		e = &Err{
			Message: v,
			Status:  http.StatusInternalServerError,
			Cause:   errors.New(v),
		}
	default:
		e = &Err{
			Message: fmt.Sprint(err),
			Status:  http.StatusInternalServerError,
		}
	}
	if e.Message == "" {
		e.Message = http.StatusText(e.Status)
	}
	return e
}

func throwErrorStatus(err error, s int, f string, args ...interface{}) {
	if err == nil {
		return
	}
	panic(Err{
		Status:  s,
		Message: fmt.Sprintf(f, args...),
		Reason:  err.Error(),
		Cause:   err,
	})
}

func throwError(err error, fmt string, args ...interface{}) {
	if err == nil {
		return
	}
	throwErrorStatus(err, detectStatusCode(err), fmt, args...)
}

func swallowError(err error, fmt string, args ...interface{}) {
	if err == nil {
		return
	}
	logrus.WithError(err).Warnf(fmt, args...)
}

func detectStatusCode(err error) int {
	if err == gorm.ErrRecordNotFound {
		return http.StatusNotFound
	}
	return http.StatusInternalServerError
}
