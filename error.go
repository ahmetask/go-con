package gocon

import (
	"fmt"
)

type TypeModel struct {
	Name string
}

var Types = map[string]TypeModel{
	"DB": {
		Name: "DB",
	},
	"Json": {
		Name: "Json",
	},
	"ConfigServer": {
		Name: "ConfigServer",
	},
	"ConfigClient": {
		Name: "ConfigClient",
	},
	"Validation": {
		Name: "Validation",
	},
	"BroadcastModule": {
		Name: "BroadcastModule",
	},
}

var EType = struct {
	DB              TypeModel
	Json            TypeModel
	ConfigServer    TypeModel
	ConfigClient    TypeModel
	Validation      TypeModel
	BroadcastModule TypeModel
}{
	DB:              Types["DB"],
	Json:            Types["Json"],
	ConfigServer:    Types["ConfigServer"],
	ConfigClient:    Types["ConfigClient"],
	Validation:      Types["Validation"],
	BroadcastModule: Types["BroadcastModule"],
}

type Error struct {
	Type  TypeModel
	Cause string
}

func (e *Error) Error() string {
	return fmt.Sprintf("Type:%s Cause:%s", e.Type.Name, e.Cause)
}

func NewError(t TypeModel, e error, messages ...string) *Error {
	msg := ""

	if e != nil {
		msg = e.Error()
	} else {
		for _, m := range messages {
			msg += m
		}
	}
	return &Error{
		Type:  t,
		Cause: msg,
	}
}
