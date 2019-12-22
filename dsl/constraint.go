package dsl

import "reflect"

type Constraint interface {
    Eval(x interface{}) bool
}

type AnyConstraint struct {
}

func (this *AnyConstraint) Eval(x interface{}) bool {
    return true
}

type EqConstraint struct {
    x interface{}
}

func (this *EqConstraint) Eval(x interface{}) bool {
    return reflect.DeepEqual(this.x, x)
}
