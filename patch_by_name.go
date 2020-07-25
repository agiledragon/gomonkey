package gomonkey

import (
	"reflect"

	"github.com/henrylee2cn/go-forceexport"
)

func ApplyFuncByTargetName(targetSymtabName string, double interface{}) *Patches {
	return create().ApplyFuncByTargetName(targetSymtabName, double)
}

func (this *Patches) ApplyFuncByTargetName(targetSymtabName string, double interface{}) *Patches {
	d := reflect.ValueOf(double)
	t := reflect.New(d.Type())
	err := forceexport.GetFunc(t.Interface(), targetSymtabName)
	if err != nil {
		panic(err)
	}
	t = reflect.ValueOf(t.Elem().Interface())
	return this.ApplyCore(t, d)
}

func ApplyFuncByDoubleName(target interface{}, doubleSymtabName string) *Patches {
	return create().ApplyFuncByDoubleName(target, doubleSymtabName)
}

func (this *Patches) ApplyFuncByDoubleName(target interface{}, doubleSymtabName string) *Patches {
	t := reflect.ValueOf(target)
	d := reflect.New(t.Type())
	err := forceexport.GetFunc(d.Interface(), doubleSymtabName)
	if err != nil {
		panic(err)
	}
	d = reflect.ValueOf(d.Elem().Interface())
	return this.ApplyCore(t, d)
}
