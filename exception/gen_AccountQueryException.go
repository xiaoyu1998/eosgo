// Code generated by gotemplate. DO NOT EDIT.

package exception

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"

	"github.com/eosspark/eos-go/log"
)

// template type Exception(PARENT,CODE,WHAT)

var AccountQueryExceptionName = reflect.TypeOf(AccountQueryException{}).Name()

type AccountQueryException struct {
	_DatabaseException
	Elog log.Messages
}

func NewAccountQueryException(parent _DatabaseException, message log.Message) *AccountQueryException {
	return &AccountQueryException{parent, log.Messages{message}}
}

func (e AccountQueryException) Code() int64 {
	return 3060002
}

func (e AccountQueryException) Name() string {
	return AccountQueryExceptionName
}

func (e AccountQueryException) What() string {
	return "Account Query Exception"
}

func (e *AccountQueryException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e AccountQueryException) GetLog() log.Messages {
	return e.Elog
}

func (e AccountQueryException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e AccountQueryException) DetailMessage() string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(int(e.Code())))
	buffer.WriteString(" ")
	buffer.WriteString(e.Name())
	buffer.WriteString(": ")
	buffer.WriteString(e.What())
	buffer.WriteString("\n")
	for _, l := range e.Elog {
		buffer.WriteString("[")
		buffer.WriteString(l.GetMessage())
		buffer.WriteString("] ")
		buffer.WriteString(l.GetContext().String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (e AccountQueryException) String() string {
	return e.DetailMessage()
}

func (e AccountQueryException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3060002,
		Name: AccountQueryExceptionName,
		What: "Account Query Exception",
	}

	return json.Marshal(except)
}

func (e AccountQueryException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*AccountQueryException):
		callback(&e)
		return true
	case func(AccountQueryException):
		callback(e)
		return true
	default:
		return false
	}
}
