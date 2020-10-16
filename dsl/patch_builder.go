package dsl

import (
	"fmt"
	"reflect"

	. "github.com/agiledragon/gomonkey/v2"
)

type FuncPara struct {
	target      interface{}
	constraints []Constraint
	behaviors   []Behavior
}

type PatchBuilder struct {
	patches  *Patches
	funcPara FuncPara
}

func NewPatchBuilder(patches *Patches) *PatchBuilder {
	funcPara := FuncPara{target: nil, constraints: make([]Constraint, 0),
		behaviors: make([]Behavior, 0)}
	return &PatchBuilder{patches: patches, funcPara: funcPara}
}

func (this *PatchBuilder) Func(target interface{}) *PatchBuilder {
	this.funcPara.target = target
	return this
}

func (this *PatchBuilder) Stubs() *PatchBuilder {
	return this
}

func (this *PatchBuilder) With(matcher ...Constraint) *PatchBuilder {
	this.funcPara.constraints = append(this.funcPara.constraints, matcher...)
	return this
}

func (this *PatchBuilder) Will(behavior Behavior) *PatchBuilder {
	this.funcPara.behaviors = append(this.funcPara.behaviors, behavior)
	return this
}

func (this *PatchBuilder) Then(behavior Behavior) *PatchBuilder {
	this.funcPara.behaviors = append(this.funcPara.behaviors, behavior)
	return this
}

func (this *PatchBuilder) End() {
	funcType := reflect.TypeOf(this.funcPara.target)
	t := reflect.ValueOf(this.funcPara.target)
	d := reflect.MakeFunc(funcType, func(inputs []reflect.Value) []reflect.Value {
		matchers := this.funcPara.constraints
		for i, input := range inputs {
			if !matchers[i].Eval(input.Interface()) {
				info := fmt.Sprintf("input paras %v is not matched", input.Interface())
				panic(info)
			}
		}
		return GetResultValues(funcType, this.funcPara.behaviors[0].Apply()[0]...)
	})
	this.patches.ApplyCore(t, d)
}
