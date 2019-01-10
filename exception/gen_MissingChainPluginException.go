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

var MissingChainPluginExceptionName = reflect.TypeOf(MissingChainPluginException{}).Name()

type MissingChainPluginException struct {
	_PluginException
	Elog log.Messages
}

func NewMissingChainPluginException(parent _PluginException, message log.Message) *MissingChainPluginException {
	return &MissingChainPluginException{parent, log.Messages{message}}
}

func (e MissingChainPluginException) Code() int64 {
	return 3110005
}

func (e MissingChainPluginException) Name() string {
	return MissingChainPluginExceptionName
}

func (e MissingChainPluginException) What() string {
	return "Missing Chain Plugin"
}

func (e *MissingChainPluginException) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e MissingChainPluginException) GetLog() log.Messages {
	return e.Elog
}

func (e MissingChainPluginException) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e MissingChainPluginException) DetailMessage() string {
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

func (e MissingChainPluginException) String() string {
	return e.DetailMessage()
}

func (e MissingChainPluginException) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3110005,
		Name: MissingChainPluginExceptionName,
		What: "Missing Chain Plugin",
	}

	return json.Marshal(except)
}

func (e MissingChainPluginException) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*MissingChainPluginException):
		callback(&e)
		return true
	case func(MissingChainPluginException):
		callback(e)
		return true
	default:
		return false
	}
}
