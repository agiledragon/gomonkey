package gomonkey

import (
	"fmt"
	"reflect"
)

type FuncPara struct {
	target    interface{}
	matchers  []Matcher
	behaviors []ReturnValue
}

type PatchBuilder struct {
	patches  *Patches
	funcPara FuncPara
}

func NewPatchBuilder(patches *Patches) *PatchBuilder {
	funcPara := FuncPara{target: nil, matchers: make([]Matcher, 0)}
	return &PatchBuilder{patches: patches, funcPara: funcPara}
}

func (this *PatchBuilder) Func(target interface{}) *PatchBuilder {
	this.funcPara.target = target
	return this
}

func (this *PatchBuilder) Stubs() *PatchBuilder {
	return this
}

func (this *PatchBuilder) With(matcher ...Matcher) *PatchBuilder {
	this.funcPara.matchers = append(this.funcPara.matchers, matcher...)
	return this
}

func (this *PatchBuilder) Will(behavior ReturnValue) *PatchBuilder {
	this.funcPara.behaviors = append(this.funcPara.behaviors, behavior)
	return this
}

func (this *PatchBuilder) End() {
	funcType := reflect.TypeOf(this.funcPara.target)
	t := reflect.ValueOf(this.funcPara.target)
	d := reflect.MakeFunc(funcType, func(inputs []reflect.Value) []reflect.Value {
		matchers := this.funcPara.matchers
		for i, input := range inputs {
			if !matchers[i].Matches(input.Interface()) {
				info := fmt.Sprintf("input paras %v is not matched", input.Interface())
				panic(info)
			}
		}
		return getResultValues(funcType, this.funcPara.behaviors[0].rets...)
	})
	this.patches.applyCore(t, d)
}

func Any() Matcher {
	return &AnyMatcher{}
}

func Eq(x interface{}) Matcher {
	return &EqMatcher{x: x}
}

func Return(x ...interface{}) ReturnValue {
	r := ReturnValue{rets: make([]interface{}, 0)}
	r.rets = append(r.rets, x...)
	return r
}

type ReturnValue struct {
	rets []interface{}
}

type Matcher interface {
	Matches(x interface{}) bool
}

type AnyMatcher struct {
}

func (this *AnyMatcher) Matches(x interface{}) bool {
	return true
}

type EqMatcher struct {
	x interface{}
}

func (this *EqMatcher) Matches(x interface{}) bool {
	return reflect.DeepEqual(this.x, x)
}
