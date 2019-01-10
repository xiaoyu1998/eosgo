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

var ActionTypeExceptionName = reflect.TypeOf(ActionTypeException{}).Name()

type ActionTypeException struct {
	_ChainException
	Elog log.Messages
}

func NewActionTypeException(parent _ChainException, message log.Message) *ActionTypeException {
	return &ActionTypeException{parent, log.Messages{message}}
}

func (e ActionTypeException) Code() int64 {
	return 3010005
}

func (e ActionTypeException) Name() string {
	return ActionTypeExceptionName
}

func (e ActionTypeException) What() string {
	return "Invalid action"
}

func (e *ActionTypeException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e ActionTypeException) GetLog() log.Messages {
	return e.Elog
}

func (e ActionTypeException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e ActionTypeException) DetailMessage() string {
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

func (e ActionTypeException) String() string {
	return e.DetailMessage()
}

func (e ActionTypeException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3010005,
		Name: ActionTypeExceptionName,
		What: "Invalid action",
	}

	return json.Marshal(except)
}

func (e ActionTypeException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*ActionTypeException):
		callback(&e)
		return true
	case func(ActionTypeException):
		callback(e)
		return true
	default:
		return false
	}
}
