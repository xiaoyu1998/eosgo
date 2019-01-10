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

var MissingPendingBlockStateName = reflect.TypeOf(MissingPendingBlockState{}).Name()

type MissingPendingBlockState struct {
	_ProducerException
	Elog log.Messages
}

func NewMissingPendingBlockState(parent _ProducerException, message log.Message) *MissingPendingBlockState {
	return &MissingPendingBlockState{parent, log.Messages{message}}
}

func (e MissingPendingBlockState) Code() int64 {
	return 3170002
}

func (e MissingPendingBlockState) Name() string {
	return MissingPendingBlockStateName
}

func (e MissingPendingBlockState) What() string {
	return "Pending block state is missing"
}

func (e *MissingPendingBlockState) AppendLog(l log.Message) {
	e.Elog = append(e.Elog, l)
}

func (e MissingPendingBlockState) GetLog() log.Messages {
	return e.Elog
}

func (e MissingPendingBlockState) TopMessage() string {
	for _, l := range e.Elog {
		if msg := l.GetMessage(); msg != "" {
			return msg
		}
	}
	return e.String()
}

func (e MissingPendingBlockState) DetailMessage() string {
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

func (e MissingPendingBlockState) String() string {
	return e.DetailMessage()
}

func (e MissingPendingBlockState) MarshalJSON() ([]byte, error) {
	type Exception struct {
		Code int64  `json:"code"`
		Name string `json:"name"`
		What string `json:"what"`
	}

	except := Exception{
		Code: 3170002,
		Name: MissingPendingBlockStateName,
		What: "Pending block state is missing",
	}

	return json.Marshal(except)
}

func (e MissingPendingBlockState) Callback(f interface{}) bool {
	switch callback := f.(type) {
	case func(*MissingPendingBlockState):
		callback(&e)
		return true
	case func(MissingPendingBlockState):
		callback(e)
		return true
	default:
		return false
	}
}
