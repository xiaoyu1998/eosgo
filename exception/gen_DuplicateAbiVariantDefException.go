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

var DuplicateAbiVariantDefExceptionName = reflect.TypeOf(DuplicateAbiVariantDefException{}).Name()

type DuplicateAbiVariantDefException struct {
	_AbiException
	Elog log.Messages
}

func NewDuplicateAbiVariantDefException(parent _AbiException, message log.Message) *DuplicateAbiVariantDefException {
	return &DuplicateAbiVariantDefException{parent, log.Messages{message}}
}

func (e DuplicateAbiVariantDefException) Code() int64 {
	return 3150015
}

func (e DuplicateAbiVariantDefException) Name() string {
	return DuplicateAbiVariantDefExceptionName
}

func (e DuplicateAbiVariantDefException) What() string {
	return "Duplicate variant definition in the ABI"
}

func (e *DuplicateAbiVariantDefException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e DuplicateAbiVariantDefException) GetLog() log.Messages {
	return e.Elog
}

func (e DuplicateAbiVariantDefException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e DuplicateAbiVariantDefException) DetailMessage() string {
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

func (e DuplicateAbiVariantDefException) String() string {
	return e.DetailMessage()
}

func (e DuplicateAbiVariantDefException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3150015,
		Name: DuplicateAbiVariantDefExceptionName,
		What: "Duplicate variant definition in the ABI",
	}

	return json.Marshal(except)
}

func (e DuplicateAbiVariantDefException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*DuplicateAbiVariantDefException):
		callback(&e)
		return true
	case func(DuplicateAbiVariantDefException):
		callback(e)
		return true
	default:
		return false
	}
}
