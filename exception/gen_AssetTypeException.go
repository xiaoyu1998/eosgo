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

var AssetTypeExceptionName = reflect.TypeOf(AssetTypeException{}).Name()

type AssetTypeException struct {
	_ChainException
	Elog log.Messages
}

func NewAssetTypeException(parent _ChainException, message log.Message) *AssetTypeException {
	return &AssetTypeException{parent, log.Messages{message}}
}

func (e AssetTypeException) Code() int64 {
	return 3010011
}

func (e AssetTypeException) Name() string {
	return AssetTypeExceptionName
}

func (e AssetTypeException) What() string {
	return "Invalid asset"
}

func (e *AssetTypeException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e AssetTypeException) GetLog() log.Messages {
	return e.Elog
}

func (e AssetTypeException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e AssetTypeException) DetailMessage() string {
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

func (e AssetTypeException) String() string {
	return e.DetailMessage()
}

func (e AssetTypeException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3010011,
		Name: AssetTypeExceptionName,
		What: "Invalid asset",
	}

	return json.Marshal(except)
}

func (e AssetTypeException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*AssetTypeException):
		callback(&e)
		return true
	case func(AssetTypeException):
		callback(e)
		return true
	default:
		return false
	}
}
