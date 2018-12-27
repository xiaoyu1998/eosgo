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

var BlockValidateExceptionName = reflect.TypeOf(BlockValidateException{}).Name()

type BlockValidateException struct {
	_BlockValidateException
	Elog log.Messages
}

func NewBlockValidateException(parent _BlockValidateException, message log.Message) *BlockValidateException {
	return &BlockValidateException{parent, log.Messages{message}}
}

func (e BlockValidateException) Code() int64 {
	return 3030000
}

func (e BlockValidateException) Name() string {
	return BlockValidateExceptionName
}

func (e BlockValidateException) What() string {
	return "Block exception"
}

func (e *BlockValidateException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e BlockValidateException) GetLog() log.Messages {
	return e.Elog
}

func (e BlockValidateException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e BlockValidateException) DetailMessage() string {
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
		buffer.WriteString("]")
		buffer.WriteString("\n")
		buffer.WriteString(l.GetContext().String())
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (e BlockValidateException) String() string {
	return e.DetailMessage()
}

func (e BlockValidateException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3030000,
		Name: BlockValidateExceptionName,
		What: "Block exception",
	}

	return json.Marshal(except)
}

func (e BlockValidateException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*BlockValidateException):
		callback(&e)
		return true
	case func(BlockValidateException):
		callback(e)
		return true
	default:
		return false
	}
}