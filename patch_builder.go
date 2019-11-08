package gomonkey

import (
    "fmt"
    "reflect"
)

type FuncPara struct {
    target    interface{}
    matchers  []Matcher
    behavior  Behavior
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

func (this *PatchBuilder) Will(behavior Behavior) *PatchBuilder {
    this.funcPara.behavior = behavior
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
        return getResultValues(funcType, this.funcPara.behavior.Apply()[0]...)
    })
    this.patches.applyCore(t, d)
}

func Any() Matcher {
    return &AnyMatcher{}
}

func Eq(x interface{}) Matcher {
    return &EqMatcher{x: x}
}

func Return(x ...interface{}) Behavior {
    r := &ReturnBehavior{rets: make([]Params, 0), params: make(Params, 0)}
    r.params = append(r.params, x...)
    return r
}

func Repeat(behavior Behavior, times int) Behavior {
    return &RepeatBehavior{rets: make([]Params, 0), behavior: behavior, times: times}
}


type Behavior interface {
    Apply() []Params
}

type ReturnBehavior struct {
    rets []Params
    params Params
}

func (this *ReturnBehavior) Apply() []Params {
    this.rets = append(this.rets, this.params)
    return this.rets
}

type RepeatBehavior struct {
    rets []Params
    behavior Behavior
    times int
}

func (this *RepeatBehavior) Apply() []Params {
    for i := 0; i < this.times; i++ {
        this.rets = append(this.rets, this.behavior.Apply()[0])
    }
    return this.rets
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
